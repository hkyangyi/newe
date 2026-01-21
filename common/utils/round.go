package utils

func RandInt(len int) int {
	return int(GetUUID()[0]) % len
}
