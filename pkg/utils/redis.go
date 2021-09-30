package utils

import "strings"

func GetRedisKey(keys ...string) string {
	return strings.Join(keys, ":")
}
