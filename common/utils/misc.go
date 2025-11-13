package utils

import (
	"strings"
	"time"
)

// RemoveSliceEmpty 去除字符串切片中的空白项（会先 TrimSpace）
func RemoveSliceEmpty(arr []string) []string {
	if len(arr) == 0 {
		return arr
	}
	res := make([]string, 0, len(arr))
	for _, s := range arr {
		v := strings.TrimSpace(s)
		if v != "" {
			res = append(res, v)
		}
	}
	return res
}

// ParseDateStringToTimestamp 使用给定 layout 解析日期字符串并返回 Unix 时间戳（秒）
// 例如 layout: "2006-01-02" 或 "2006-01-02 15:04:05"
func ParseDateStringToTimestamp(dateStr, layout string) (int64, error) {
	if layout == "" {
		layout = "2006-01-02 15:04:05"
	}
	t, err := time.ParseInLocation(layout, strings.TrimSpace(dateStr), time.Local)
	if err != nil {
		return 0, err
	}
	return t.Unix(), nil
}
