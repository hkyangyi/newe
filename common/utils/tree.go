package utils

import (
	"fmt"
	"reflect"
	"strings"
)

// ListToTree 将任意列表（struct 或 map）转换为树结构。
// - list: 列表，元素可以是 struct、*struct、map[string]any
// - idField/pidField/childrenField: 字段名或 json tag 名（不区分大小写）。
// 返回值是 []map[string]any，便于 JSON 序列化和动态 children 注入。
func ListToTree(list any, idField, pidField, childrenField string) []map[string]any {
	// 预处理字段名，统一小写
	idKey := strings.ToLower(idField)
	pidKey := strings.ToLower(pidField)
	childKey := childrenField // children 字段名保持原样输出

	v := reflect.ValueOf(list)
	if v.Kind() != reflect.Slice {
		return nil
	}

	// 先把每一项转成 map[string]any，并建立 id -> item 的索引
	items := make([]map[string]any, 0, v.Len())
	index := make(map[string]map[string]any, v.Len())

	for i := 0; i < v.Len(); i++ {
		itemVal := indirect(v.Index(i))
		m := toMap(itemVal)
		// 归一化键名查找 id、pid
		id := toString(findValueByKey(m, idKey))
		// 初始化 children
		if _, ok := m[childKey]; !ok {
			m[childKey] = make([]map[string]any, 0)
		}
		items = append(items, m)
		if id != "" {
			index[id] = m
		}
		// 把 pid 也放回去（统一输出键名）
		if _, ok := m[pidField]; !ok {
			// 若原始键名不同于传入 pidField，则补充 pidField 键
			if pv, exists := findKeyOriginal(m, pidKey); exists && pidField != pv {
				m[pidField] = m[pv]
			}
		}
		if _, ok := m[idField]; !ok {
			if iv, exists := findKeyOriginal(m, idKey); exists && idField != iv {
				m[idField] = m[iv]
			}
		}
	}

	// 组装树
	roots := make([]map[string]any, 0)
	for _, m := range items {
		id := toString(findValueByKey(m, strings.ToLower(idField)))
		pid := toString(findValueByKey(m, strings.ToLower(pidField)))
		if pid == "" || index[pid] == nil || pid == id { // 根节点（无父或父不存在/自指向）
			roots = append(roots, m)
			continue
		}
		parent := index[pid]
		// 追加进父的 children
		if ch, ok := parent[childKey]; ok {
			if arr, ok2 := ch.([]map[string]any); ok2 {
				parent[childKey] = append(arr, m)
			}
		} else {
			parent[childKey] = []map[string]any{m}
		}
	}
	return roots
}

// toMap 将 struct 或 map 转换为 map[string]any。
// 对于 struct，使用字段名与 `json` tag（优先 tag）。
func toMap(v reflect.Value) map[string]any {
	m := make(map[string]any)
	if !v.IsValid() {
		return m
	}
	switch v.Kind() {
	case reflect.Map:
		iter := v.MapRange()
		for iter.Next() {
			k := iter.Key()
			if k.Kind() == reflect.String {
				m[k.String()] = toInterface(iter.Value())
			}
		}
	case reflect.Struct:
		t := v.Type()
		for i := 0; i < v.NumField(); i++ {
			f := t.Field(i)
			// 跳过未导出字段
			if f.PkgPath != "" {
				continue
			}
			key := f.Name
			if tag := f.Tag.Get("json"); tag != "" && tag != "-" {
				// 取 tag 的第一个逗号前子串
				if idx := strings.IndexByte(tag, ','); idx >= 0 {
					key = tag[:idx]
				} else {
					key = tag
				}
			}
			m[key] = toInterface(v.Field(i))
		}
	default:
		// 其他类型直接返回空 map
	}
	return m
}

func toInterface(v reflect.Value) any {
	v = indirect(v)
	if !v.IsValid() {
		return nil
	}
	switch v.Kind() {
	case reflect.Struct:
		// 对嵌套 struct 递归为 map
		return toMap(v)
	case reflect.Slice, reflect.Array:
		arr := make([]any, v.Len())
		for i := 0; i < v.Len(); i++ {
			arr[i] = toInterface(v.Index(i))
		}
		return arr
	case reflect.Map:
		return toMap(v)
	default:
		return v.Interface()
	}
}

func indirect(v reflect.Value) reflect.Value {
	for v.IsValid() && v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return reflect.Value{}
		}
		v = v.Elem()
	}
	return v
}

// findValueByKey 在 map 中以不区分大小写的方式查找键
func findValueByKey(m map[string]any, lowerKey string) any {
	for k, v := range m {
		if strings.ToLower(k) == lowerKey {
			return v
		}
	}
	return nil
}

// findKeyOriginal 返回与小写匹配的原始键名
func findKeyOriginal(m map[string]any, lowerKey string) (string, bool) {
	for k := range m {
		if strings.ToLower(k) == lowerKey {
			return k, true
		}
	}
	return "", false
}

func toString(v any) string {
	switch x := v.(type) {
	case string:
		return x
	case []byte:
		return string(x)
	case nil:
		return ""
	default:
		return strings.TrimSpace(fmt.Sprint(v))
	}
}
