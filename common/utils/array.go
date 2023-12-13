package utils

import (
	"strconv"
	"strings"
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
