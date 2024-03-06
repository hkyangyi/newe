package utils

import (
	"bytes"
	"log"
	"strconv"
	"strings"
	"unicode"
)

func InArrar(key string, arr []string) bool {

	for _, v := range arr {
		if key == v {
			return true
		}
	}

	return false
}

func BoolToInt(b bool) int {
	if b {
		return 1
	} else {
		return -1
	}
}

func JoinInt(ints []int, sep string) string {
	var arrs []string
	for _, v := range ints {
		is := strconv.Itoa(v)
		arrs = append(arrs, is)
	}
	res := strings.Join(arrs, sep)
	return res
}

func SplitInt(str, sep string) []int {
	arrs := strings.Split(str, sep)
	var ints []int
	for _, v := range arrs {
		ii, _ := strconv.Atoi(v)
		ints = append(ints, ii)
	}
	return ints
}

func JoinInt64(ints []int64, sep string) string {
	var arrs []string
	for _, v := range ints {
		is := strconv.Itoa(int(v))
		arrs = append(arrs, is)
	}
	res := strings.Join(arrs, sep)
	return res
}

func SplitInt64(str, sep string) []int64 {
	arrs := strings.Split(str, sep)
	var ints []int64
	for _, v := range arrs {
		ii, _ := strconv.ParseInt(v, 10, 64)
		ints = append(ints, ii)
	}
	return ints
}

// 驼峰式写法转为下划线写法
func Camel2Case(name string) string {
	buffer := NewBuffer()
	for i, r := range name {
		if unicode.IsUpper(r) {
			if i != 0 {
				buffer.Append('_')
			}
			buffer.Append(unicode.ToLower(r))
		} else {
			buffer.Append(r)
		}
	}
	return buffer.String()
}

// 首字母大写
func Ucfirst(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+1:]
	}
	return ""
}

// 首字母小写
func Lcfirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}

// 内嵌bytes.Buffer，支持连写
type Buffer struct {
	*bytes.Buffer
}

func NewBuffer() *Buffer {
	return &Buffer{Buffer: new(bytes.Buffer)}
}

func (b *Buffer) Append(i interface{}) *Buffer {
	switch val := i.(type) {
	case int:
		b.append(strconv.Itoa(val))
	case int64:
		b.append(strconv.FormatInt(val, 10))
	case uint:
		b.append(strconv.FormatUint(uint64(val), 10))
	case uint64:
		b.append(strconv.FormatUint(val, 10))
	case string:
		b.append(val)
	case []byte:
		b.Write(val)
	case rune:
		b.WriteRune(val)
	}
	return b
}

func (b *Buffer) append(s string) *Buffer {
	defer func() {
		if err := recover(); err != nil {
			log.Println("*****内存不够了！******")
		}
	}()
	b.WriteString(s)
	return b
}
