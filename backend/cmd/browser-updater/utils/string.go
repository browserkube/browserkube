package utils

import (
	"encoding/base64"
	"strings"
)

func ParseStringMap(s, separator string) map[string]string {
	ss := strings.Split(s, separator)
	m := make(map[string]string, len(ss))
	for _, v := range ss {
		z := strings.Split(v, "=")

		// remove any quotes from the string
		m[strings.ReplaceAll(z[0], "\"", "")] = strings.ReplaceAll(z[1], "\"", "")
	}
	return m
}

func BasicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func RemoveFirst(text, separator string) string {
	splits := strings.Split(text, separator)
	if len(splits) == 1 {
		return text
	}
	var str string
	for i := 1; i < len(splits); i++ {
		str += splits[i]
	}
	return str
}
