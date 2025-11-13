package model

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hkyangyi/newe/common/utils"
)

// ParseCreateTable 解析 SHOW CREATE TABLE 返回的创建语句，返回表级信息、列列表和主键列名
func ParseCreateTable(createSQL string) (DevDbtable, []DevDbtableColumns, []string, error) {
	var table DevDbtable

	// 表级信息
	table.TableRemark = parseCommentFromSQL(createSQL)
	table.Engine = parseEngineFromSQL(createSQL)
	table.Charset = parseCharsetFromSQL(createSQL)
	table.CreateTime = time.Now().Unix()

	// 列解析
	l := strings.Index(createSQL, "(")
	r := strings.LastIndex(strings.ToUpper(createSQL), ") ENGINE")
	if r == -1 {
		r = strings.LastIndex(createSQL, ")")
	}
	body := createSQL
	if l >= 0 && r > l {
		body = createSQL[l+1 : r]
	}

	lines := splitSQLLinesGeneric(body)

	// 主键
	pkCols := findPrimaryKeysFromSQL(createSQL)

	cols := make([]DevDbtableColumns, 0, len(lines))
	colLineRe := regexp.MustCompile(`^\s*` + "`" + `([^` + "`" + `]*)` + "`" + `\s+(.+)$`)
	commentRe := regexp.MustCompile("COMMENT '([^']*)'")
	typeRe := regexp.MustCompile(`^([a-zA-Z0-9_]+)(\(([^)]*)\))?`)

	sortIdx := 0
	for _, ln := range lines {
		ln = strings.TrimSpace(ln)
		if ln == "" {
			continue
		}
		up := strings.ToUpper(ln)
		if strings.HasPrefix(up, "PRIMARY KEY") || strings.HasPrefix(up, "UNIQUE KEY") || strings.HasPrefix(up, "KEY ") || strings.HasPrefix(up, "CONSTRAINT") {
			continue
		}
		m := colLineRe.FindStringSubmatch(ln)
		if len(m) == 0 {
			continue
		}
		name := m[1]
		rest := strings.TrimSpace(m[2])

		colType := ""
		colLen := 0
		colScale := 0
		isUnsigned := 0
		notNull := 0
		if mt := typeRe.FindStringSubmatch(rest); len(mt) > 0 {
			colType = strings.ToLower(mt[1])
			if len(mt) > 3 && mt[3] != "" {
				parts := strings.Split(mt[3], ",")
				if n, err := strconv.Atoi(strings.TrimSpace(parts[0])); err == nil {
					colLen = n
				}
				if len(parts) > 1 {
					if s, err := strconv.Atoi(strings.TrimSpace(parts[1])); err == nil {
						colScale = s
					}
				}
			}
		}

		def := ""
		defRe := regexp.MustCompile(`DEFAULT\s+((?:'[^']*')|[^\s]+)`)
		if md := defRe.FindStringSubmatch(rest); len(md) > 1 {
			val := strings.Trim(md[1], "'")
			if strings.ToUpper(val) == "NULL" {
				def = ""
			} else {
				def = val
			}
		}

		extra := ""
		if strings.Contains(strings.ToLower(rest), "auto_increment") {
			extra = "auto_increment"
		}
		if strings.Contains(strings.ToLower(rest), "unsigned") {
			isUnsigned = 1
		}
		if strings.Contains(strings.ToUpper(rest), "NOT NULL") {
			notNull = 1
		}

		comment := ""
		if mc := commentRe.FindStringSubmatch(rest); len(mc) > 1 {
			comment = mc[1]
		}

		col := DevDbtableColumns{
			Id:           utils.GetUUID(),
			TableId:      "",
			Column:       name,
			ColumnType:   colType,
			ColumnLong:   colLen,
			ColumnScale:  colScale,
			IsUnsigned:   isUnsigned,
			NotNull:      notNull,
			ColumnValue:  def,
			ColumnKey:    "",
			Extra:        extra,
			ColumnCom:    comment,
			Sort:         sortIdx,
			IsSearch:     0,
			IsDict:       0,
			DictCode:     "",
			FormComplate: "",
			FormProps:    "",
			CreateTime:   time.Now().Unix(),
		}
		cols = append(cols, col)
		sortIdx++
	}

	// 标记主键
	for i := range cols {
		for _, pk := range pkCols {
			if cols[i].Column == pk {
				cols[i].ColumnKey = "PRI"
			}
		}
	}

	return table, cols, pkCols, nil
}

// BuildCreateTableSQL 根据 DevDbtable 和 DevDbtableColumns 生成 CREATE TABLE SQL（简单实现）
func BuildCreateTableSQL(table DevDbtable, cols []DevDbtableColumns) (string, error) {
	if table.DbTableName == "" {
		return "", fmt.Errorf("table name empty")
	}
	lines := make([]string, 0, len(cols))
	pkCols := make([]string, 0)
	for _, c := range cols {
		typ := c.ColumnType
		if c.ColumnScale > 0 {
			typ = fmt.Sprintf("%s(%d,%d)", typ, c.ColumnLong, c.ColumnScale)
		} else if c.ColumnLong > 0 {
			typ = fmt.Sprintf("%s(%d)", typ, c.ColumnLong)
		}
		if c.IsUnsigned == 1 {
			typ += " unsigned"
		}
		nullStr := ""
		if c.NotNull == 1 {
			nullStr = " NOT NULL"
		}
		defStr := ""
		if c.ColumnValue != "" {
			// 判断是否数字简单处理
			if _, err := strconv.ParseFloat(c.ColumnValue, 64); err == nil {
				defStr = " DEFAULT " + c.ColumnValue
			} else {
				defStr = " DEFAULT '" + strings.ReplaceAll(c.ColumnValue, "'", "\\'") + "'"
			}
		}
		extra := ""
		if c.Extra != "" {
			extra = " " + c.Extra
		}
		comment := ""
		if c.ColumnCom != "" {
			comment = " COMMENT '" + strings.ReplaceAll(c.ColumnCom, "'", "\\'") + "'"
		}
		line := fmt.Sprintf("`%s` %s%s%s%s%s", c.Column, typ, nullStr, defStr, extra, comment)
		lines = append(lines, line)
		if c.ColumnKey == "PRI" {
			pkCols = append(pkCols, "`"+c.Column+"`")
		}
	}
	if len(pkCols) > 0 {
		lines = append(lines, "PRIMARY KEY ("+strings.Join(pkCols, ",")+")")
	}
	ddl := "CREATE TABLE `" + table.DbTableName + "` (\n" + strings.Join(lines, ",\n") + "\n)"
	if table.Engine != "" {
		ddl += " ENGINE=" + table.Engine
	}
	if table.Charset != "" {
		ddl += " DEFAULT CHARSET=" + table.Charset
	}
	if table.TableRemark != "" {
		ddl += " COMMENT='" + strings.ReplaceAll(table.TableRemark, "'", "\\'") + "'"
	}
	return ddl, nil
}

// 内部帮助函数
func parseEngineFromSQL(sql string) string {
	re := regexp.MustCompile(`ENGINE=([A-Za-z0-9_]+)`)
	if m := re.FindStringSubmatch(sql); len(m) > 1 {
		return m[1]
	}
	return ""
}

func parseCharsetFromSQL(sql string) string {
	re := regexp.MustCompile(`DEFAULT CHARSET=([A-Za-z0-9_]+)`)
	if m := re.FindStringSubmatch(sql); len(m) > 1 {
		return m[1]
	}
	return ""
}

func parseCommentFromSQL(sql string) string {
	re := regexp.MustCompile("COMMENT='([^']*)'")
	if m := re.FindStringSubmatch(sql); len(m) > 1 {
		return m[1]
	}
	return ""
}

func findPrimaryKeysFromSQL(sql string) []string {
	re := regexp.MustCompile(`PRIMARY KEY\s*\(([^)]*)\)`)
	if m := re.FindStringSubmatch(sql); len(m) > 1 {
		parts := strings.Split(m[1], ",")
		res := make([]string, 0, len(parts))
		for _, p := range parts {
			p = strings.TrimSpace(p)
			p = strings.Trim(p, "` ")
			if p != "" {
				res = append(res, p)
			}
		}
		return res
	}
	return nil
}

func splitSQLLinesGeneric(body string) []string {
	res := []string{}
	var cur strings.Builder
	depth := 0
	inQuotes := false
	var quoteChar rune
	for _, ch := range body {
		if inQuotes {
			cur.WriteRune(ch)
			if ch == quoteChar {
				inQuotes = false
			}
			continue
		}
		if ch == '\'' || ch == '"' || ch == '`' {
			inQuotes = true
			quoteChar = ch
			cur.WriteRune(ch)
			continue
		}
		if ch == '(' {
			depth++
			cur.WriteRune(ch)
			continue
		}
		if ch == ')' {
			if depth > 0 {
				depth--
			}
			cur.WriteRune(ch)
			continue
		}
		if ch == ',' && depth == 0 {
			res = append(res, cur.String())
			cur.Reset()
			continue
		}
		cur.WriteRune(ch)
	}
	tail := strings.TrimSpace(cur.String())
	if tail != "" {
		res = append(res, tail)
	}
	return res
}
