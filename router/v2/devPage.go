package v2

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/hkyangyi/newe/common/utils"
	"github.com/hkyangyi/newe/model"
	"github.com/hkyangyi/newe/router/app"
)

func DevPageRegRoute(r *gin.RouterGroup) {
	r.POST("add", DevPageAdd)
	r.PUT("edit", DevPageEdit)
	r.DELETE("del", DevPageDel)
	r.GET("page", DevPageGetPage)
	r.GET("list", DevPageGetList)
	r.POST("generate", DevPageGenerate)
}

// 添加
func DevPageAdd(c *gin.Context) {
	var g = app.NewApp(c)
	var a model.DevPage
	if err := g.Bind(&a); err != nil {
		g.Error(err)
		return
	}
	if err := a.Add(); err != nil {
		g.Error(err)
		return
	}
	g.SUCCESS(nil)
}

// 编辑
func DevPageEdit(c *gin.Context) {
	var g = app.NewApp(c)
	var a model.DevPage
	if err := g.Bind(&a); err != nil {
		g.Error(err)
		return
	}
	if err := a.Edit(); err != nil {
		g.Error(err)
		return
	}
	g.SUCCESS(nil)
}

// 删除
func DevPageDel(c *gin.Context) {
	var g = app.NewApp(c)
	var a model.DevPage
	if err := g.Bind(&a); err != nil {
		g.Error(err)
		return
	}
	if err := a.Del(); err != nil {
		g.Error(err)
		return
	}
	g.SUCCESS(nil)
}

// 获取分页列表
func DevPageGetPage(c *gin.Context) {
	var g = app.NewApp(c)
	var a model.DevPage
	if err := g.Bind(&a); err != nil {
		g.Error(err)
		return
	}

	var params []interface{}
	var where []string

	page := utils.PageList{
		Page:     a.Page,
		PageSize: a.PageSize,
	}
	ws := strings.Join(where, " AND ")
	res := a.GetPage(page, ws, params...)
	g.SUCCESS(res)
}

// 获取所有列表
func DevPageGetList(c *gin.Context) {
	var g = app.NewApp(c)
	var a model.DevPage
	if err := g.Bind(&a); err != nil {
		g.Error(err)
		return
	}
	res := a.GetList()
	g.SUCCESS(res)
}

// 生成页面
func DevPageGenerate(c *gin.Context) {
	var g = app.NewApp(c)
	var req model.DevPage
	if err := g.Bind(&req); err != nil {
		g.Error(err)
		return
	}
	if err := req.Refresh(); err != nil {
		g.Error(err)
		return
	}
	// 获取字段列表：从 DevDbtable 表中取列
	var tb = model.DevDbtable{ID: req.SqlTable}
	if err := tb.RefResh(); err != nil {
		g.Error(err)
		return
	}
	columns := tb.GetColumns()
	if req.PageType != 3 {
		if err := DevPageAcGenerate(req, tb, columns); err != nil {
			g.Error(err)
			return
		}

		if err := DevPageAcGenerateFrontendAPI(req, tb, columns); err != nil {
			g.Error(err)
			return
		}

		g.SUCCESS(nil)
	} else {
		if err := DevPageAcGenerateOneToOne(req, tb, columns); err != nil {
			g.Error(err)
			return
		}
		g.SUCCESS(nil)
	}
}

// 生成页面：生成后端 model 与前端简单页面模板，文件保存到 devgen/ 目录
func DevPageAcGenerate(pd model.DevPage, tb model.DevDbtable, columns []model.DevDbtableColumns) error {
	// 目录
	base := "assets/devgen/" + tb.DbTableName
	backendDir := filepath.Join(base, "go")
	if err := os.MkdirAll(backendDir, 0755); err != nil {
		return err
	}
	modelName := utils.Ucfirst(snakeToCamel(tb.DbTableName))
	// 生成后端 model 文件（使用模板）
	modelFile := filepath.Join(backendDir, "model_"+modelName+".go")
	// build struct string
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("type %s struct {\n", modelName))
	// ensure Id field
	sb.WriteString("    Id string `json:\"id\" form:\"id\" gorm:\"primary_key\"`\n")
	for _, c := range columns {
		if c.Column == "id" {
			continue
		}
		fieldName := utils.Ucfirst(snakeToCamel(c.Column))
		// default to string type; include json tag and comment
		sb.WriteString(fmt.Sprintf("    %s %s `json:\"%s\" form:\"%s\"` // %s\n", fieldName, dbtypetogo(c.ColumnType), utils.Lcfirst(snakeToCamel(c.Column)), utils.Lcfirst(snakeToCamel(c.Column)), strings.ReplaceAll(c.ColumnCom, "`", "'")))
	}
	// add timestamps
	sb.WriteString("    utils.PageList \n")
	sb.WriteString("}\n")

	tpl, err := template.New("model").Parse(DevPageModelTpl)
	if err != nil {
		return err
	}
	data := map[string]interface{}{
		"ModelName": modelName,
		"TableName": pd.SqlTable,
		"StStr":     sb.String(),
	}
	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil {
		return err
	}
	if err := os.WriteFile(modelFile, buf.Bytes(), 0644); err != nil {
		return err
	}

	//生成CORS
	corsFile := filepath.Join(backendDir, "cors_"+modelName+".go")
	corsTpl, err := template.New("cors").Parse(DevPageCorsTpl)
	if err != nil {
		return err
	}

	//查询条件
	var wherestr strings.Builder
	for _, c := range columns {
		if c.IsSearch == 1 {

			switch c.ColumnType {
			case "varchar", "text", "char":
				wherestr.WriteString(fmt.Sprintf("        if a.%s != \"\" {\n", utils.Ucfirst(snakeToCamel(c.Column))))
				wherestr.WriteString(fmt.Sprintf("            where = append(where, \"%s LIKE ?\")\n", c.Column))
				wherestr.WriteString(fmt.Sprintf("            params = append(params, \"%%\"+a.%s+\"%%\")\n", utils.Ucfirst(snakeToCamel(c.Column))))
				wherestr.WriteString("        }\n")
			case "int", "bigint", "tinyint", "smallint", "mediumint", "float", "double", "decimal":
				wherestr.WriteString(fmt.Sprintf("        if a.%s != 0 {\n", utils.Ucfirst(snakeToCamel(c.Column))))
				wherestr.WriteString(fmt.Sprintf("            where = append(where, \"%s = ?\")\n", c.Column))
				wherestr.WriteString(fmt.Sprintf("            params = append(params, a.%s)\n", utils.Ucfirst(snakeToCamel(c.Column))))
				wherestr.WriteString("        }\n")
			default:
				wherestr.WriteString(fmt.Sprintf("        if a.%s != \"\" {\n", utils.Ucfirst(snakeToCamel(c.Column))))
				wherestr.WriteString(fmt.Sprintf("            where = append(where, \"%s = ?\")\n", c.Column))
				wherestr.WriteString(fmt.Sprintf("            params = append(params, a.%s)\n", utils.Ucfirst(snakeToCamel(c.Column))))
				wherestr.WriteString("        }\n")
			}
		}
	}

	data["WhereStr"] = wherestr.String()

	var corsBuf bytes.Buffer

	if err := corsTpl.Execute(&corsBuf, data); err != nil {
		return err
	}
	if err := os.WriteFile(corsFile, corsBuf.Bytes(), 0644); err != nil {
		return err
	}
	// 生成前端 data.ts
	if err := DevPageAcGenerateFrontendData(pd, tb, columns); err != nil {
		return err
	}

	return nil
}

// 生成前端 data.ts（表格/表单/搜索定义）
func DevPageAcGenerateFrontendData(pd model.DevPage, tb model.DevDbtable, columns []model.DevDbtableColumns) error {
	// 目录
	base := "assets/devgen/" + tb.DbTableName
	frontendDir := filepath.Join(base, "vue")
	if err := os.MkdirAll(frontendDir, 0755); err != nil {
		return err
	}

	dataFile := filepath.Join(frontendDir, "data.ts")
	df, err := os.Create(dataFile)
	if err != nil {
		return err
	}
	defer df.Close()
	modelName := utils.Ucfirst(snakeToCamel(tb.DbTableName))

	// 使用 pd 填充基本信息

	fmt.Fprintln(df, "import type { VbenFormSchema } from '#/adapter/form';")
	fmt.Fprintln(df, "import type { VxeTableGridOptions } from '#/adapter/vxe-table';")
	fmt.Fprintln(df)
	fmt.Fprintln(df, "import { formatDate } from '#/plugins/utils';")
	fmt.Fprintln(df, "import { GetTree } from './api';")
	fmt.Fprintln(df, "// 结构体")
	fmt.Fprintf(df, "export interface %s {\n", modelName)
	for _, c := range columns {
		fieldName := utils.Lcfirst(snakeToCamel(c.Column))
		fileType := gotypetots(c.ColumnType)
		fmt.Fprintf(df, "  %s: %s; // %s\n", fieldName, fileType, strings.ReplaceAll(c.ColumnCom, "`", "'"))
	}
	fmt.Fprintln(df, "}")
	fmt.Fprintln(df)
	fmt.Fprintf(df, "export const %sColumns: VxeTableGridOptions['columns'] = [", modelName)
	for k, c := range columns {
		fieldName := utils.Lcfirst(snakeToCamel(c.Column))
		//判断是否包含time
		if strings.Contains(c.Column, "time") {
			fmt.Fprintf(df, "  { field: '%s', title: '%s', width: 180, formatter: ({ cellValue }) => formatDate(cellValue), },\n", fieldName, strings.ReplaceAll(c.ColumnCom, "`", "'"))
			continue
		}
		if k == 0 && pd.PageType == 2 {
			fmt.Fprintf(df, "  { field: '%s', title: '%s',  treeNode: true, minWidth: 150, align: 'center' },\n", fieldName, strings.ReplaceAll(c.ColumnCom, "`", "'"))
			continue
		}
		if c.FormComplate == "Switch" {
			fmt.Fprintf(df, "  { field: '%s', title: '%s', minWidth: 150, align: 'center', slots: { default: '%s' }},\n", fieldName, strings.ReplaceAll(c.ColumnCom, "`", "'"), fieldName)
			continue
		}
		fmt.Fprintf(df, "  { field: '%s', title: '%s', minWidth: 150, align: 'center' },\n", fieldName, strings.ReplaceAll(c.ColumnCom, "`", "'"))
	}

	fmt.Fprintln(df, " { field: 'operation', title: '操作', minWidth: 150, align: 'center', slots: { default: 'operation' }, },\n")

	fmt.Fprintln(df, "];")

	// readFormSchemas
	fmt.Fprintf(df, "export const %sFormSchemas: VbenFormSchema[] = [", modelName)
	for _, c := range columns {
		fieldName := utils.Lcfirst(snakeToCamel(c.Column))
		if c.Column == "id" {
			fmt.Fprintln(df, "  { fieldName: 'id', label: 'ID', component: 'Input', dependencies: { show: false, triggerFields: ['id'] }, },")
			continue
		}
		if c.IsFormshow != 1 {
			fmt.Fprintf(df, "//")
		}
		switch c.FormComplate {
		case "InputNumber":
			fmt.Fprintf(df, "  { fieldName: '%s', component: 'InputNumber', label: '%s', componentProps: { placeholder: '请输入%s' }, },\n", fieldName, strings.ReplaceAll(c.ColumnCom, "`", "'"), strings.ReplaceAll(c.ColumnCom, "`", "'"))
			continue
		case "DatePicker":
			fmt.Fprintf(df, "  { fieldName: '%s', component: 'DatePicker', label: '%s', componentProps: { placeholder: '请输入%s', valueFormat: 'YYYY-MM-DD HH:mm:ss' }, },\n", fieldName, strings.ReplaceAll(c.ColumnCom, "`", "'"), strings.ReplaceAll(c.ColumnCom, "`", "'"))
			continue
		case "Switch":
			fmt.Fprintf(df, "  { fieldName: '%s', component: 'Switch', label: '%s', componentProps: { placeholder: '请输入%s',checkedValue: 1, unCheckedValue: -1  }, },\n", fieldName, strings.ReplaceAll(c.ColumnCom, "`", "'"), strings.ReplaceAll(c.ColumnCom, "`", "'"))
			continue
		case "Select":
			fmt.Fprintf(df, "  { fieldName: '%s', component: 'Select', label: '%s', componentProps: { placeholder: '请输入%s' , options: [ { label: 'varchar', value: 'varchar' },] }, },\n", fieldName, strings.ReplaceAll(c.ColumnCom, "`", "'"), strings.ReplaceAll(c.ColumnCom, "`", "'"))
			continue
		case "TextArea":
			fmt.Fprintf(df, "  { fieldName: '%s', component: 'InputTextArea', label: '%s', componentProps: { placeholder: '请输入%s', rows: 4 }, },\n", fieldName, strings.ReplaceAll(c.ColumnCom, "`", "'"), strings.ReplaceAll(c.ColumnCom, "`", "'"))
			continue
		case "RadioGroup":
			fmt.Fprintf(df, "  { fieldName: '%s', component: 'RadioGroup', label: '%s', componentProps: { placeholder: '请输入%s' , options: [ { label: 'varchar', value: 'varchar' },] }, },\n", fieldName, strings.ReplaceAll(c.ColumnCom, "`", "'"), strings.ReplaceAll(c.ColumnCom, "`", "'"))
			continue
		case "CheckboxGroup":
			fmt.Fprintf(df, "  { fieldName: '%s', component: 'CheckboxGroup', label: '%s', componentProps: { placeholder: '请输入%s' , options: [ { label: 'varchar', value: 'varchar' },] }, },\n", fieldName, strings.ReplaceAll(c.ColumnCom, "`", "'"), strings.ReplaceAll(c.ColumnCom, "`", "'"))
			continue
		case "TimePicker":
			fmt.Fprintf(df, "  { fieldName: '%s', component: 'TimePicker', label: '%s', componentProps: { placeholder: '请输入%s', valueFormat: 'HH:mm:ss' }, },\n", fieldName, strings.ReplaceAll(c.ColumnCom, "`", "'"), strings.ReplaceAll(c.ColumnCom, "`", "'"))
			continue
		case "ColorPicker":
			fmt.Fprintf(df, "  { fieldName: '%s', component: 'ColorPicker', label: '%s', componentProps: { placeholder: '请输入%s' }, },\n", fieldName, strings.ReplaceAll(c.ColumnCom, "`", "'"), strings.ReplaceAll(c.ColumnCom, "`", "'"))
			continue
		case "Rate":
			fmt.Fprintf(df, "  { fieldName: '%s', component: 'Rate', label: '%s', componentProps: { placeholder: '请输入%s' }, },\n", fieldName, strings.ReplaceAll(c.ColumnCom, "`", "'"), strings.ReplaceAll(c.ColumnCom, "`", "'"))
			continue
		case "ApiTreeSelect":
			fmt.Fprintf(df, "  { fieldName: '%s', component: 'ApiTreeSelect', label: '%s', componentProps: { placeholder: '请输入%s', api: GetTree, fieldNames: { label: 'name', value: 'id', children: 'children' }, style: { width: %s }, }, },\n", fieldName, strings.ReplaceAll(c.ColumnCom, "`", "'"), strings.ReplaceAll(c.ColumnCom, "`", "'"), "'100%'")
			continue
		default:
			fmt.Fprintf(df, "  { fieldName: '%s', component: '%s', label: '%s', componentProps: { placeholder: '请输入%s' }, },\n", fieldName, c.FormComplate, strings.ReplaceAll(c.ColumnCom, "`", "'"), strings.ReplaceAll(c.ColumnCom, "`", "'"))
		}

	}
	fmt.Fprintln(df, "];")

	// SearchSchema
	fmt.Fprintln(df, "export const SearchSchema: VbenFormSchema[] = [")
	for _, c := range columns {
		if c.IsSearch != 1 {
			continue
		}
		fieldName := utils.Lcfirst(snakeToCamel(c.Column))
		switch c.ColumnType {
		case "varchar", "text", "char":
			fmt.Fprintf(df, "  { fieldName: '%s', component: 'Input', label: '%s', componentProps: { placeholder: '按%s搜索' }, },\n", fieldName, strings.ReplaceAll(c.ColumnCom, "`", "'"), strings.ReplaceAll(c.ColumnCom, "`", "'"))
		case "int", "bigint", "tinyint", "smallint", "mediumint", "float", "double", "decimal":
			fmt.Fprintf(df, "  { fieldName: '%s', component: 'InputNumber', label: '%s', componentProps: { placeholder: '按%s搜索' }, },\n", fieldName, strings.ReplaceAll(c.ColumnCom, "`", "'"), strings.ReplaceAll(c.ColumnCom, "`", "'"))
		default:
			fmt.Fprintf(df, "  { fieldName: '%s', component: 'Input', label: '%s', componentProps: { placeholder: '按%s搜索' }, },\n", fieldName, strings.ReplaceAll(c.ColumnCom, "`", "'"), strings.ReplaceAll(c.ColumnCom, "`", "'"))
		}

	}
	fmt.Fprintln(df, "];")

	return nil
}

// 生成前端apits
func DevPageAcGenerateFrontendAPI(pd model.DevPage, tb model.DevDbtable, columns []model.DevDbtableColumns) error {
	// 目录
	base := "assets/devgen/" + tb.DbTableName
	frontendDir := filepath.Join(base, "vue")
	if err := os.MkdirAll(frontendDir, 0755); err != nil {
		return err
	}

	// 生成前端 API 文件（使用模板）
	apiFile := filepath.Join(frontendDir, "api.ts")
	apiTpl, err := template.New("api").Parse(DevPageVueApiTpl)
	if err != nil {
		return err
	}
	modelName := utils.Ucfirst(snakeToCamel(tb.DbTableName))

	//rootpath 前面有/后面没有/ 进行检测
	rootpath := pd.RouteRoot
	if !strings.HasPrefix(rootpath, "/") {
		rootpath = "/" + rootpath
	}
	rootpath = strings.TrimSuffix(rootpath, "/")

	data := map[string]interface{}{
		"ModelName": modelName,
		"RootPath":  rootpath,
	}
	var apiBuf bytes.Buffer
	if err := apiTpl.Execute(&apiBuf, data); err != nil {
		return err
	}
	if err := os.WriteFile(apiFile, apiBuf.Bytes(), 0644); err != nil {
		return err
	}

	// 生成 index.vue
	indexFile := filepath.Join(frontendDir, "index.vue")
	idxF, err := os.Create(indexFile)
	if err != nil {
		return err
	}
	defer idxF.Close()

	// 导入 api 并别名化函数为 <ModelName>GetPage / Edit / Del 风格
	apiAlias := modelName
	// 写入文件内容
	fmt.Fprintf(idxF, "<script lang=\"ts\" setup>\n")
	fmt.Fprintf(idxF, "import type { %s as Record } from './data';\n", modelName)
	fmt.Fprintf(idxF, "\n")
	fmt.Fprintf(idxF, "import type { VxeTableGridOptions } from '#/adapter/vxe-table';\n")
	fmt.Fprintf(idxF, "\n")
	fmt.Fprintf(idxF, "import { Page, useVbenModal } from '@vben/common-ui';\n")
	fmt.Fprintf(idxF, "import { Plus } from '@vben/icons';\n\n")
	fmt.Fprintf(idxF, "import { Button, Switch } from 'ant-design-vue';\n\n")
	fmt.Fprintf(idxF, "import { useVbenVxeGrid } from '#/adapter/vxe-table';\n\n")
	if pd.PageType == 2 {
		fmt.Fprintf(idxF, "import { Add as %sAdd, Edit as %sEdit, GetTree as %sGetTree, Del as %sDel } from './api';\n", apiAlias, apiAlias, apiAlias, apiAlias)
	} else {
		fmt.Fprintf(idxF, "import { Add as %sAdd, Edit as %sEdit, GetPage as %sGetPage, Del as %sDel } from './api';\n", apiAlias, apiAlias, apiAlias, apiAlias)
	}

	fmt.Fprintf(idxF, "import { SearchSchema, %sColumns } from './data';\n", modelName)
	fmt.Fprintf(idxF, "import FormModal from './formModal.vue';\n\n")
	fmt.Fprintf(idxF, "defineOptions({ name: '%s' });\n\n", utils.Ucfirst(snakeToCamel(tb.DbTableName)))

	fmt.Fprintf(idxF, "const [Modal, modalApi] = useVbenModal({\n")
	fmt.Fprintf(idxF, "  connectedComponent: FormModal,\n")
	fmt.Fprintf(idxF, "  destroyOnClose: true,\n")
	fmt.Fprintf(idxF, "});\n\n")

	fmt.Fprintf(idxF, "const [Grid, gridApi] = useVbenVxeGrid({\n")
	fmt.Fprintf(idxF, "  formOptions: {\n")
	fmt.Fprintf(idxF, "    collapsed: false,\n")
	fmt.Fprintf(idxF, "    schema: SearchSchema,\n")
	fmt.Fprintf(idxF, "    showCollapseButton: true,\n")
	fmt.Fprintf(idxF, "    submitButtonOptions: { content: '查询' },\n")
	fmt.Fprintf(idxF, "    submitOnChange: false,\n")
	fmt.Fprintf(idxF, "    submitOnEnter: true,\n")
	fmt.Fprintf(idxF, "  },\n")
	if pd.PageType == 2 {
		fmt.Fprintf(idxF, "  gridOptions: {\n")
		fmt.Fprintf(idxF, "    columns: %sColumns as any,\n", modelName)
		fmt.Fprintf(idxF, "    height: 'auto',\n")
		fmt.Fprintf(idxF, "    pagerConfig: { enabled: false, pageSize: 20, pageSizes: [10, 20, 50, 100] },\n")
		fmt.Fprintf(idxF, "    rowConfig: { keyField: 'id' },\n")
		fmt.Fprintf(idxF, "    toolbarConfig: { custom: true, export: false, refresh: true, zoom: true },\n")
		fmt.Fprintf(idxF, "    cellConfig: { height: 80 },\n")
		fmt.Fprintf(idxF, "    proxyConfig: { ajax: { query: async ( formValues ) => {  const items = await %sGetTree({ ...formValues }); return { items }; } } },\n", apiAlias)
		fmt.Fprintf(idxF, "  treeConfig: {\n")
		fmt.Fprintf(idxF, "    parentField: 'pid',\n")
		fmt.Fprintf(idxF, "    rowField: 'id',\n")
		fmt.Fprintf(idxF, "    transform: false,\n")
		fmt.Fprintf(idxF, "    showLine: true,\n")
		fmt.Fprintf(idxF, "    expandAll: true,\n")
		fmt.Fprintf(idxF, "  },\n")
		fmt.Fprintf(idxF, "} as VxeTableGridOptions,\n")
		fmt.Fprintf(idxF, "});\n\n")
	} else {
		fmt.Fprintf(idxF, "  gridOptions: {\n")
		fmt.Fprintf(idxF, "    columns: %sColumns as any,\n", modelName)
		fmt.Fprintf(idxF, "    height: 'auto',\n")
		fmt.Fprintf(idxF, "    pagerConfig: { enabled: true, pageSize: 20, pageSizes: [10, 20, 50, 100] },\n")
		fmt.Fprintf(idxF, "    rowConfig: { keyField: 'id' },\n")
		fmt.Fprintf(idxF, "    toolbarConfig: { custom: true, export: false, refresh: true, zoom: true },\n")
		fmt.Fprintf(idxF, "    cellConfig: { height: 80 },\n")
		fmt.Fprintf(idxF, "    proxyConfig: { ajax: { query: async ( { page: pageInfo = { currentPage: 1, pageSize: 10 } }, formValues ) => { const { currentPage, pageSize } = pageInfo; const res = await %sGetPage({ ...formValues, page: currentPage, pageSize }); return { items: res.list, total: res.total, pageSize: res.pageSize, currentPage: res.page }; } } },\n", apiAlias)
		fmt.Fprintf(idxF, "  } as VxeTableGridOptions,\n")
		fmt.Fprintf(idxF, "});\n\n")
	}

	fmt.Fprintf(idxF, "function onRefresh() { gridApi.query(); }\n")
	fmt.Fprintf(idxF, "function onCreate() { modalApi.setData({}).open(); }\n\n")
	fmt.Fprintf(idxF, "function onEdit(row: Record) { modalApi.setData({ ...row }).open(); }\n\n")
	fmt.Fprintf(idxF, "async function onDelete(row: Record) { await %sDel(row); onRefresh(); }\n\n", apiAlias)
	fmt.Fprintf(idxF, "async function onToggle(row: Record, checked: any) { const val = checked ? 1 : -1; await %sEdit({ id: row.id, status: val }); row.status = val; }\n", apiAlias)

	fmt.Fprintf(idxF, "</script>\n")
	fmt.Fprintf(idxF, "<template>\n  <Page auto-content-height>\n    <Grid>\n      <template #toolbar-tools>\n        <Button type=\"primary\" @click=\"onCreate()\">\n          <Plus class=\"size-5\" /> 新增\n        </Button>\n      </template>\n      <template #status=\"row\">\n        <Switch :checked=\"row.row.status === 1\" checked-children=\"启用\" un-checked-children=\"禁用\" @change=\"(val) => onToggle(row.row, val)\" />\n      </template>\n      <template #operation=\"{ row }\">\n        <Button type=\"link\" size=\"small\" @click=\"onEdit(row)\">编辑</Button>\n        <Button type=\"link\" size=\"small\" danger @click=\"onDelete(row)\">删除</Button>\n      </template>\n    </Grid>\n    <Modal @success=\"onRefresh\" />\n  </Page>\n</template>\n")

	// 生成 formModal.vue
	formFile := filepath.Join(frontendDir, "formModal.vue")
	fmF, err := os.Create(formFile)
	if err != nil {
		return err
	}
	defer fmF.Close()

	fmt.Fprintf(fmF, "<script lang=\"ts\" setup>\n")
	fmt.Fprintf(fmF, "import { ref } from 'vue';\n\n")
	fmt.Fprintf(fmF, "import { useVbenModal } from '@vben/common-ui';\n\n")
	fmt.Fprintf(fmF, "import { useVbenForm } from '#/adapter/form';\n\n")
	fmt.Fprintf(fmF, "import { Add as %sAdd, Edit as %sEdit } from './api';\n", apiAlias, apiAlias)
	fmt.Fprintf(fmF, "import { %sFormSchemas } from './data';\n\n", modelName)

	fmt.Fprintf(fmF, "const emit = defineEmits<{ success: [] }>();\n\n")
	fmt.Fprintf(fmF, "const isUpdate = ref(false);\n\n")
	fmt.Fprintf(fmF, "const [Modal, modalApi] = useVbenModal({\n")
	fmt.Fprintf(fmF, "  onConfirm: () => onSubmit(),\n")
	fmt.Fprintf(fmF, "  onCancel: () => modalApi.close(),\n")
	fmt.Fprintf(fmF, "  onOpened() {\n")
	fmt.Fprintf(fmF, "    // 每次打开时重置表单\n")
	fmt.Fprintf(fmF, "    formApi.resetForm();\n")
	fmt.Fprintf(fmF, "    const origin = modalApi.getData<any>() || {};\n")
	fmt.Fprintf(fmF, "    isUpdate.value = Boolean(origin?.id);\n")
	fmt.Fprintf(fmF, "    formApi.setValues(isUpdate.value ? origin : { pid: origin.pid || '' });\n")
	fmt.Fprintf(fmF, "  },\n")
	fmt.Fprintf(fmF, "});\n\n")

	fmt.Fprintf(fmF, "const [Form, formApi] = useVbenForm({\n")
	fmt.Fprintf(fmF, "  schema: %sFormSchemas,\n", utils.Ucfirst(modelName))
	fmt.Fprintf(fmF, "  showDefaultActions: false,\n")
	fmt.Fprintf(fmF, "  commonConfig: { colon: true },\n")
	fmt.Fprintf(fmF, "});\n\n")

	fmt.Fprintf(fmF, "async function onSubmit() {\n")
	fmt.Fprintf(fmF, "  await formApi.validate();\n")
	fmt.Fprintf(fmF, "  const values = await formApi.getValues();\n")
	fmt.Fprintf(fmF, "  await (isUpdate.value ? %sEdit(values) : %sAdd(values));\n", apiAlias, apiAlias)
	fmt.Fprintf(fmF, "  modalApi.close();\n")
	fmt.Fprintf(fmF, "  emit('success');\n")
	fmt.Fprintf(fmF, "}\n")
	fmt.Fprintf(fmF, "</script>\n")
	fmt.Fprintf(fmF, "<template>\n  <Modal :title=\"isUpdate ? '编辑' : '新建'\" class=\"w-full max-w-[720px]\">\n    <Form class=\"mx-4\" />\n  </Modal>\n</template>\n")

	return nil
}

func snakeToCamel(s string) string {
	parts := strings.Split(s, "_")
	var b strings.Builder
	for _, p := range parts {
		if p == "" {
			continue
		}
		b.WriteString(strings.ToUpper(p[:1]))
		if len(p) > 1 {
			b.WriteString(p[1:])
		}
	}
	return b.String()
}

func gotypetots(t string) string {
	switch t {
	case "int", "bigint", "tinyint", "smallint", "mediumint":
		return "number"
	case "float", "double", "decimal":
		return "number"
	case "datetime", "timestamp", "date":
		return "number"
	default:
		return "string"
	}
}

// 一对一模板
func DevPageAcGenerateOneToOne(pd model.DevPage, tb model.DevDbtable, columns []model.DevDbtableColumns) error {
	// 目录
	// 目录
	base := "assets/devgen/" + tb.DbTableName
	backendDir := filepath.Join(base, "go")
	if err := os.MkdirAll(backendDir, 0755); err != nil {
		return err
	}
	modelName := utils.Ucfirst(snakeToCamel(tb.DbTableName))
	// 生成后端 model 文件（使用模板）
	modelFile := filepath.Join(backendDir, "model_"+modelName+".go")
	// build struct string
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("type %s struct {\n", modelName))
	// ensure Id field
	sb.WriteString("    Id string `json:\"id\" form:\"id\" gorm:\"primary_key\"`\n")
	for _, c := range columns {
		if c.Column == "id" {
			continue
		}
		fieldName := utils.Ucfirst(snakeToCamel(c.Column))
		// default to string type; include json tag and comment
		sb.WriteString(fmt.Sprintf("    %s %s `json:\"%s\" form:\"%s\"` // %s\n", fieldName, dbtypetogo(c.ColumnType), utils.Lcfirst(snakeToCamel(c.Column)), utils.Lcfirst(snakeToCamel(c.Column)), strings.ReplaceAll(c.ColumnCom, "`", "'")))
	}
	sb.WriteString("    utils.PageList \n")
	sb.WriteString("}\n")

	tpl, err := template.New("model").Parse(DevPageModelTpl)
	if err != nil {
		return err
	}
	data := map[string]interface{}{
		"ModelName": modelName,
		"TableName": pd.SqlTable,
		"StStr":     sb.String(),
	}
	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil {
		return err
	}
	if err := os.WriteFile(modelFile, buf.Bytes(), 0644); err != nil {
		return err
	}

	//生成CORS
	corsFile := filepath.Join(backendDir, "cors_"+modelName+".go")
	corsTpl, err := template.New("cors").Parse(DevPageCorsTpl)
	if err != nil {
		return err
	}

	//查询条件
	var wherestr strings.Builder
	for _, c := range columns {
		if c.IsSearch == 1 {

			switch c.ColumnType {
			case "varchar", "text", "char":
				wherestr.WriteString(fmt.Sprintf("        if a.%s != \"\" {\n", utils.Ucfirst(snakeToCamel(c.Column))))
				wherestr.WriteString(fmt.Sprintf("            where = append(where, \"%s LIKE ?\")\n", c.Column))
				wherestr.WriteString(fmt.Sprintf("            params = append(params, \"%%\"+a.%s+\"%%\")\n", utils.Ucfirst(snakeToCamel(c.Column))))
				wherestr.WriteString("        }\n")
			case "int", "bigint", "tinyint", "smallint", "mediumint", "float", "double", "decimal":
				wherestr.WriteString(fmt.Sprintf("        if a.%s != 0 {\n", utils.Ucfirst(snakeToCamel(c.Column))))
				wherestr.WriteString(fmt.Sprintf("            where = append(where, \"%s = ?\")\n", c.Column))
				wherestr.WriteString(fmt.Sprintf("            params = append(params, a.%s)\n", utils.Ucfirst(snakeToCamel(c.Column))))
				wherestr.WriteString("        }\n")
			default:
				wherestr.WriteString(fmt.Sprintf("        if a.%s != \"\" {\n", utils.Ucfirst(snakeToCamel(c.Column))))
				wherestr.WriteString(fmt.Sprintf("            where = append(where, \"%s = ?\")\n", c.Column))
				wherestr.WriteString(fmt.Sprintf("            params = append(params, a.%s)\n", utils.Ucfirst(snakeToCamel(c.Column))))
				wherestr.WriteString("        }\n")
			}
		}
	}

	data["WhereStr"] = wherestr.String()

	var corsBuf bytes.Buffer

	if err := corsTpl.Execute(&corsBuf, data); err != nil {
		return err
	}
	if err := os.WriteFile(corsFile, corsBuf.Bytes(), 0644); err != nil {
		return err
	}

	//生成API.TS
	if err := DevPageAcGenerateFrontendAPI(pd, tb, columns); err != nil {
		return err
	}

	//生成data.ts
	if err := DevPageAcGenerateFrontendData(pd, tb, columns); err != nil {
		return err
	}

	//生成index.vue
	if err := DevPageAcGenerateFrontendOneToOneVue(pd, tb, columns); err != nil {
		return err
	}
	//生成formModal.vue
	if err := DevPageAcGenerateFrontendOneToOneFormModalVue(pd, tb, columns); err != nil {
		return err
	}
	//生成leftTree.vue
	if err := DevPageAcGenerateFrontendLeftTreeVue(pd, tb, columns); err != nil {
		return err
	}

	return nil
}

// 生成前端一对一index.vue
func DevPageAcGenerateFrontendOneToOneVue(pd model.DevPage, tb model.DevDbtable, columns []model.DevDbtableColumns) error {
	// 目录
	base := "assets/devgen/" + tb.DbTableName
	frontendDir := filepath.Join(base, "vue")
	if err := os.MkdirAll(frontendDir, 0755); err != nil {
		return err
	}
	// 生成 index.vue
	indexFile := filepath.Join(frontendDir, "index.vue")
	idxF, err := os.Create(indexFile)
	if err != nil {
		return err
	}
	defer idxF.Close()
	modelName := utils.Ucfirst(snakeToCamel(tb.DbTableName))
	// 导入 api 并别名化函数为 <ModelName>GetPage / Edit / Del 风格
	apiAlias := modelName
	// 写入文件内容

	fmt.Fprintf(idxF, "<script lang=\"ts\" setup>\n")
	fmt.Fprintf(idxF, "import type { %s as Record } from './data';\n", apiAlias)
	fmt.Fprintf(idxF, "\n")
	fmt.Fprintf(idxF, "import type { VxeTableGridOptions } from '#/adapter/vxe-table';\n")
	fmt.Fprintf(idxF, "\n")

	fmt.Fprintf(idxF, "import { ref } from 'vue';")
	fmt.Fprintf(idxF, "\n")

	fmt.Fprintf(idxF, "import { Page, useVbenDrawer, useVbenModal } from '@vben/common-ui';")
	fmt.Fprintf(idxF, "\n")
	fmt.Fprintf(idxF, "import { Plus } from '@vben/icons';")
	fmt.Fprintf(idxF, "\n")
	fmt.Fprintf(idxF, "import { Button, Card } from 'ant-design-vue';\n")
	fmt.Fprintf(idxF, "import { useVbenVxeGrid } from '#/adapter/vxe-table';\n")
	fmt.Fprintf(idxF, "import { Del, GetPage } from './api';\n")
	fmt.Fprintf(idxF, "import FormModal from './formModal.vue';")
	fmt.Fprintf(idxF, "\n")
	fmt.Fprintf(idxF, "import leftTree from './leftTree.vue';")
	fmt.Fprintf(idxF, "\n")

	fmt.Fprintf(idxF, "import { %sColumns, SearchSchema } from './data';\n", apiAlias)
	fmt.Fprintf(idxF, "\n")
	fmt.Fprintf(idxF, "defineOptions({ name: '%s' });\n", apiAlias)

	fmt.Fprintf(idxF, "const classId = ref<string | undefined>(undefined);\n")
	fmt.Fprintf(idxF, "\n")
	fmt.Fprintf(idxF, "const [Modal, modalApi] = useVbenModal({\n")
	fmt.Fprintf(idxF, "connectedComponent: FormModal,\n")
	fmt.Fprintf(idxF, "destroyOnClose: true,\n")
	fmt.Fprintf(idxF, "});\n")
	fmt.Fprintf(idxF, "\n")

	fmt.Fprintf(idxF, "const [Grid, gridApi] = useVbenVxeGrid({\n")
	fmt.Fprintf(idxF, "  formOptions: {\n")
	fmt.Fprintf(idxF, "    collapsed: false,\n")
	fmt.Fprintf(idxF, "    schema: SearchSchema,\n")
	fmt.Fprintf(idxF, "   // 控制表单是否显示折叠按钮\n")
	fmt.Fprintf(idxF, "    showCollapseButton: true,\n")
	fmt.Fprintf(idxF, "    submitButtonOptions: {\n")
	fmt.Fprintf(idxF, "      content: '查询',\n")
	fmt.Fprintf(idxF, "    },\n")
	fmt.Fprintf(idxF, "    // 是否在字段值改变时提交表单\n")
	fmt.Fprintf(idxF, "    submitOnChange: false,\n")
	fmt.Fprintf(idxF, "   // 按下回车时是否提交表单\n")
	fmt.Fprintf(idxF, "    submitOnEnter: true,\n")
	fmt.Fprintf(idxF, "  },\n")

	fmt.Fprintf(idxF, "  gridOptions: {\n")
	fmt.Fprintf(idxF, "    columns: %sColumns as any,\n", apiAlias)
	fmt.Fprintf(idxF, "    height: 'auto',\n")
	fmt.Fprintf(idxF, "    pagerConfig: { enabled: true, pageSize: 10, pageSizes: [10, 20, 50, 100] },\n")
	fmt.Fprintf(idxF, "    rowConfig: { keyField: 'id' },\n")
	fmt.Fprintf(idxF, "    toolbarConfig: { custom: true, export: false, refresh: true, zoom: true },\n")
	fmt.Fprintf(idxF, "    cellConfig: {\n")
	fmt.Fprintf(idxF, "      height: 80,\n")
	fmt.Fprintf(idxF, "    },\n")
	fmt.Fprintf(idxF, "\n")
	fmt.Fprintf(idxF, "    proxyConfig: {\n")
	fmt.Fprintf(idxF, "      ajax: {\n")
	fmt.Fprintf(idxF, "        query: async ({ page: pageInfo }, formValues) => {\n")
	fmt.Fprintf(idxF, "          const { currentPage, pageSize } = pageInfo;\n")
	fmt.Fprintf(idxF, "          const res = await GetPage({\n")
	fmt.Fprintf(idxF, "            ...formValues,\n")
	fmt.Fprintf(idxF, "            page: currentPage,\n")
	fmt.Fprintf(idxF, "            pageSize,\n")
	fmt.Fprintf(idxF, "          });\n")
	fmt.Fprintf(idxF, "          return {\n")
	fmt.Fprintf(idxF, "            items: res.list,\n")
	fmt.Fprintf(idxF, "            total: res.total,\n")
	fmt.Fprintf(idxF, "            pageSize: res.pageSize,\n")
	fmt.Fprintf(idxF, "            currentPage: res.page,\n")
	fmt.Fprintf(idxF, "          };\n")
	fmt.Fprintf(idxF, "        },\n")
	fmt.Fprintf(idxF, "      },\n")
	fmt.Fprintf(idxF, "    },\n")
	fmt.Fprintf(idxF, "  } as VxeTableGridOptions,\n")
	fmt.Fprintf(idxF, "});\n")
	fmt.Fprintf(idxF, "\n")
	fmt.Fprintf(idxF, "function onRefresh() {\n")
	fmt.Fprintf(idxF, "  gridApi.query();\n")
	fmt.Fprintf(idxF, "}\n")
	fmt.Fprintf(idxF, "function onCreate() {\n")
	fmt.Fprintf(idxF, "  modalApi.setData({ %s:classId.value }).open();\n", utils.Lcfirst(snakeToCamel(pd.SqlKey)))
	fmt.Fprintf(idxF, "}\n")
	fmt.Fprintf(idxF, "\n")
	fmt.Fprintf(idxF, "function onEdit(row: Record) {\n")
	fmt.Fprintf(idxF, "  modalApi.setData({ ...row }).open();\n")
	fmt.Fprintf(idxF, "}\n")
	fmt.Fprintf(idxF, "\n")
	fmt.Fprintf(idxF, "async function onDelete(row: Record) {\n")
	fmt.Fprintf(idxF, "  await Del(row);\n")
	fmt.Fprintf(idxF, "  onRefresh();\n")
	fmt.Fprintf(idxF, "}\n")
	fmt.Fprintf(idxF, "\n")
	fmt.Fprintf(idxF, "async function handleTreeChange(id?: string) {\n")
	fmt.Fprintf(idxF, "  classId.value = id || '';\n")
	fmt.Fprintf(idxF, "  await gridApi.formApi.setValues({ %s: id || '' });\n", utils.Lcfirst(snakeToCamel(pd.SqlKey)))
	fmt.Fprintf(idxF, "  gridApi.formApi.submitForm();\n")
	fmt.Fprintf(idxF, "}\n")
	fmt.Fprintf(idxF, "\n")
	fmt.Fprintf(idxF, "</script>\n")
	fmt.Fprintf(idxF, "<template>\n")
	fmt.Fprintf(idxF, "  <Page auto-content-height>\n")
	fmt.Fprintf(idxF, "    <div class=\"flex h-full flex-row gap-3\">\n")
	fmt.Fprintf(idxF, "      <Card class=\"w-1/6 flex-1\">\n")
	fmt.Fprintf(idxF, "        <leftTree\n")
	fmt.Fprintf(idxF, "          @change=\"handleTreeChange\"\n")
	fmt.Fprintf(idxF, "          @refresh=\"() => (classId = undefined)\"\n")
	fmt.Fprintf(idxF, "        />\n")
	fmt.Fprintf(idxF, "      </Card>\n")
	fmt.Fprintf(idxF, "      <Grid class=\"h-full w-5/6 flex-1\">\n")
	fmt.Fprintf(idxF, "        <template #toolbar-tools>\n")
	fmt.Fprintf(idxF, "          <Button type=\"primary\" @click=\"onCreate()\">\n")
	fmt.Fprintf(idxF, "            <Plus class=\"size-5\" /> 添加\n")
	fmt.Fprintf(idxF, "          </Button>\n")
	fmt.Fprintf(idxF, "        </template>\n")

	for _, c := range columns {
		if c.Column == pd.SqlKey || c.Column == "id" {
			continue
		}
		if c.FormComplate == "Switch" {
			fmt.Fprintf(idxF, "        <template #%s=\"row\">\n", utils.Lcfirst(snakeToCamel(c.Column)))
			fmt.Fprintf(idxF, "          <Switch :checked=\"row.row.%s === 1\" checked-children=\"是\" un-checked-children=\"否\" disabled />\n", utils.Lcfirst(snakeToCamel(c.Column)))
			fmt.Fprintf(idxF, "        </template>\n")
		}
	}

	fmt.Fprintf(idxF, "        <template #operation=\"{ row }\">\n")
	fmt.Fprintf(idxF, "          <Button type=\"link\" size=\"small\" @click=\"onEdit(row)\"> 编辑 </Button>\n")
	fmt.Fprintf(idxF, "          <Button type=\"link\" size=\"small\" danger @click=\"onDelete(row)\">\n")
	fmt.Fprintf(idxF, "            删除\n")
	fmt.Fprintf(idxF, "          </Button>\n")
	fmt.Fprintf(idxF, "        </template>\n")
	fmt.Fprintf(idxF, "      </Grid>\n")
	fmt.Fprintf(idxF, "      <Modal @success=\"onRefresh\" />\n")
	fmt.Fprintf(idxF, "    </div>\n")
	fmt.Fprintf(idxF, "  </Page>\n")
	fmt.Fprintf(idxF, "</template>\n")
	fmt.Fprintf(idxF, "\n")
	fmt.Fprintf(idxF, "<style scoped></style>\n")

	return nil
}

// 生成leftTree.vue
func DevPageAcGenerateFrontendLeftTreeVue(pd model.DevPage, tb model.DevDbtable, columns []model.DevDbtableColumns) error {
	// 目录
	base := "assets/devgen/" + tb.DbTableName
	frontendDir := filepath.Join(base, "vue")
	if err := os.MkdirAll(frontendDir, 0755); err != nil {
		return err
	}
	// 生成 index.vue
	indexFile := filepath.Join(frontendDir, "leftTree.vue")
	idxF, err := os.Create(indexFile)
	if err != nil {
		return err
	}
	defer idxF.Close()
	//查看关联的页面是否存在

	// treetb:=model.DevPage{SqlTable: pd.SqlTableDeputy}
	// treetb.GetOneByTable()

	var tbd = model.DevDbtable{ID: pd.SqlTableDeputy}
	if err := tbd.RefResh(); err != nil {
		return err
	}
	modelNamed := utils.Ucfirst(snakeToCamel(tbd.DbTableName))
	// 导入 api 并别名化函数为 <ModelName>GetPage / Edit / Del 风格
	// 写入文件内容
	fmt.Fprintf(idxF, "<script lang=\"ts\" setup>\n")

	fmt.Fprintln(idxF, "import type { TreeProps } from 'ant-design-vue';")
	fmt.Fprintf(idxF, "\n")
	fmt.Fprintln(idxF, "import { onMounted, ref } from 'vue';")
	fmt.Fprintf(idxF, "\n")
	fmt.Fprintln(idxF, "import { Button, Empty, Tree } from 'ant-design-vue';")
	fmt.Fprintf(idxF, "\n")
	fmt.Fprintf(idxF, "import { GetTree } from './../%s/api';\n", modelNamed)
	fmt.Fprintf(idxF, "\n")
	fmt.Fprintln(idxF, " interface Emits {")
	fmt.Fprintln(idxF, "   (e: 'change', deptId?: string): void;")
	fmt.Fprintln(idxF, "   (e: 'refresh'): void;")
	fmt.Fprintln(idxF, " }")
	fmt.Fprintf(idxF, "\n")
	fmt.Fprintln(idxF, "const emit = defineEmits<Emits>();")
	fmt.Fprintf(idxF, "\n")
	fmt.Fprintln(idxF, "const loading = ref(false);")
	fmt.Fprintf(idxF, "\n")
	fmt.Fprintln(idxF, "const treeData = ref<any[]>([]);")
	fmt.Fprintf(idxF, "\n")
	fmt.Fprintln(idxF, "const selectedKeys = ref<string[]>([]);")
	fmt.Fprintf(idxF, "\n")

	fmt.Fprintln(idxF, "async function loadTree() {")
	fmt.Fprintf(idxF, "\n")
	fmt.Fprintln(idxF, "  loading.value = true;")
	fmt.Fprintf(idxF, "\n")
	fmt.Fprintln(idxF, "  try {")
	fmt.Fprintf(idxF, "\n")
	fmt.Fprintln(idxF, "    const list = await GetTree();")
	fmt.Fprintf(idxF, "\n")
	fmt.Fprintln(idxF, "    treeData.value = Array.isArray(list) ? list : [];")
	fmt.Fprintf(idxF, "\n")
	fmt.Fprintln(idxF, "  } finally {")
	fmt.Fprintf(idxF, "\n")
	fmt.Fprintln(idxF, "    loading.value = false;")
	fmt.Fprintf(idxF, "\n")
	fmt.Fprintln(idxF, "  }")
	fmt.Fprintf(idxF, "\n")
	fmt.Fprintln(idxF, "}")
	fmt.Fprintf(idxF, "\n")
	fmt.Fprintln(idxF)
	fmt.Fprintln(idxF, "function onSelect(keys: any[]) {")
	fmt.Fprintf(idxF, "\n")
	fmt.Fprintln(idxF, "  selectedKeys.value = keys as string[];")
	fmt.Fprintf(idxF, "\n")
	fmt.Fprintln(idxF, "  emit('change', selectedKeys.value[0]);")
	fmt.Fprintf(idxF, "\n")
	fmt.Fprintln(idxF, "}")
	fmt.Fprintf(idxF, "\n")
	fmt.Fprintln(idxF, "function onRefresh() {")
	fmt.Fprintf(idxF, "\n")
	fmt.Fprintln(idxF, "  emit('refresh');")
	fmt.Fprintf(idxF, "\n")
	fmt.Fprintln(idxF, "}")
	fmt.Fprintf(idxF, "\n")
	fmt.Fprintln(idxF, "onMounted(loadTree);")
	fmt.Fprintf(idxF, "\n")
	fmt.Fprintln(idxF, "const fieldNames: TreeProps['fieldNames'] = {")
	fmt.Fprintf(idxF, "\n")
	fmt.Fprintln(idxF, "  key: 'id',")
	fmt.Fprintf(idxF, "\n")
	fmt.Fprintln(idxF, "  title: 'name',")
	fmt.Fprintf(idxF, "\n")
	fmt.Fprintln(idxF, "  children: 'children',")
	fmt.Fprintf(idxF, "\n")
	fmt.Fprintln(idxF, "};")
	fmt.Fprintf(idxF, "\n")
	fmt.Fprintf(idxF, "</script>\n")

	fmt.Fprintf(idxF, "<template>\n")
	fmt.Fprintf(idxF, "  <div class=\"flex h-full flex-col gap-2\">\n")
	fmt.Fprintf(idxF, "    <div class=\"flex items-center justify-between\">\n")
	fmt.Fprintf(idxF, "      <div class=\"font-medium\">类型</div>\n")
	fmt.Fprintf(idxF, "      <Button size=\"small\" @click=\"onRefresh\">刷新</Button>\n")
	fmt.Fprintf(idxF, "    </div>\n")
	fmt.Fprintf(idxF, "    <div class=\"flex-1 overflow-auto\">\n")
	fmt.Fprintf(idxF, "      <Tree\n")
	fmt.Fprintf(idxF, "        v-if=\"treeData.length > 0\"\n")
	fmt.Fprintf(idxF, "        v-model:selected-keys=\"selectedKeys\"\n")

	fmt.Fprintf(idxF, "        :tree-data=\"treeData\"\n")
	fmt.Fprintf(idxF, "        :field-names=\"fieldNames\"\n")
	fmt.Fprintf(idxF, "        :loading=\"loading\"\n")
	fmt.Fprintf(idxF, "        show-line\n")
	fmt.Fprintf(idxF, "        default-expand-all\n")
	fmt.Fprintf(idxF, "        @select=\"onSelect\"\n")
	fmt.Fprintf(idxF, "      />\n")
	fmt.Fprintf(idxF, "      <Empty v-else />\n")
	fmt.Fprintf(idxF, "    </div>\n")
	fmt.Fprintf(idxF, "  </div>\n")
	fmt.Fprintf(idxF, "</template>\n")
	fmt.Fprintf(idxF, "<style scoped></style>)\n")

	return nil
}

// 生成FormModal.vue
func DevPageAcGenerateFrontendOneToOneFormModalVue(pd model.DevPage, tb model.DevDbtable, columns []model.DevDbtableColumns) error {
	// 目录
	base := "assets/devgen/" + tb.DbTableName
	frontendDir := filepath.Join(base, "vue")
	if err := os.MkdirAll(frontendDir, 0755); err != nil {
		return err
	}
	// 生成 formModal.vue
	formFile := filepath.Join(frontendDir, "formModal.vue")
	fmF, err := os.Create(formFile)
	if err != nil {
		return err
	}
	defer fmF.Close()
	modelName := utils.Ucfirst(snakeToCamel(tb.DbTableName))

	fmt.Fprintf(fmF, "<script lang=\"ts\" setup>\n")
	fmt.Fprintf(fmF, "import { ref } from 'vue';\n\n")
	fmt.Fprintf(fmF, "import { useVbenModal } from '@vben/common-ui';\n\n")
	fmt.Fprintf(fmF, "import { useVbenForm } from '#/adapter/form';\n\n")
	fmt.Fprintf(fmF, "import { Add as %sAdd, Edit as %sEdit } from './api';\n", modelName, modelName)
	fmt.Fprintf(fmF, "import { %sFormSchemas } from './data';\n\n", modelName)

	fmt.Fprintf(fmF, "const emit = defineEmits<{ success: [] }>();\n\n")
	fmt.Fprintf(fmF, "const isUpdate = ref(false);\n\n")
	fmt.Fprintf(fmF, "const [Modal, modalApi] = useVbenModal({\n")
	fmt.Fprintf(fmF, "  onConfirm: () => onSubmit(),\n")
	fmt.Fprintf(fmF, "  onCancel: () => modalApi.close(),\n")
	fmt.Fprintf(fmF, "  onOpened() {\n")
	fmt.Fprintf(fmF, "    // 每次打开时重置表单\n")
	fmt.Fprintf(fmF, "    formApi.resetForm();\n")
	fmt.Fprintf(fmF, "    const origin = modalApi.getData<any>() || {};\n")
	fmt.Fprintf(fmF, "    isUpdate.value = Boolean(origin?.id);\n")
	fmt.Fprintf(fmF, "    formApi.setValues(isUpdate.value ? origin : { pid: origin.pid || '' });\n")
	fmt.Fprintf(fmF, "  },\n")
	fmt.Fprintf(fmF, "});\n\n")

	fmt.Fprintf(fmF, "const [Form, formApi] = useVbenForm({\n")
	fmt.Fprintf(fmF, "  schema: %sFormSchemas,\n", utils.Ucfirst(modelName))
	fmt.Fprintf(fmF, "  showDefaultActions: false,\n")
	fmt.Fprintf(fmF, "  commonConfig: { colon: true },\n")
	fmt.Fprintf(fmF, "});\n\n")

	fmt.Fprintf(fmF, "async function onSubmit() {\n")
	fmt.Fprintf(fmF, "  await formApi.validate();\n")
	fmt.Fprintf(fmF, "  const values = await formApi.getValues();\n")
	fmt.Fprintf(fmF, "  await (isUpdate.value ? %sEdit(values) : %sAdd(values));\n", modelName, modelName)
	fmt.Fprintf(fmF, "  modalApi.close();\n")
	fmt.Fprintf(fmF, "  emit('success');\n")
	fmt.Fprintf(fmF, "}\n")
	fmt.Fprintf(fmF, "</script>\n")
	fmt.Fprintf(fmF, "<template>\n  <Modal :title=\"isUpdate ? '编辑' : '新建'\" class=\"w-full max-w-[720px]\">\n    <Form class=\"mx-4\" />\n  </Modal>\n</template>\n")

	return nil
}

// 数据库类型转GO类型
func dbtypetogo(t string) string {
	switch t {
	case "int", "tinyint", "smallint", "mediumint":
		return "int"
	case "bigint":
		return "int64"
	case "float", "double", "decimal":
		return "float64"
	case "datetime", "timestamp", "date":
		return "time.Time"
	default:
		return "string"
	}
}
