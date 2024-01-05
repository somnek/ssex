package main

import (
	"regexp"
	"strconv"
)

func StrToInt(s string) (int, error) {
	return strconv.Atoi(s)
}

func IsNumber(s string) bool {
	re := regexp.MustCompile(`^[0-9]+$`)
	return re.MatchString(s)
}
