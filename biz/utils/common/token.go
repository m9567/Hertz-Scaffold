package common

import (
	"strings"
)

var (
	split = "__"
)

func GetUserNameByToken(token string) string {
	return strings.Split(token, split)[0][6:]
}
