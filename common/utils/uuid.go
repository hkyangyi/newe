package utils

import (
	"strings"

	"github.com/google/uuid"
)

func GetUUID() string {
	u := uuid.New().String()
	uu := strings.Replace(u, "-", "", -1)
	return uu
}
