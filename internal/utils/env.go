package utils

import (
	"os"
	"strconv"
	"strings"
)

func GetEnvStr(key string) (string, error) {
	v := os.Getenv(key)
	return v, nil
}

func GetEnvInt(key string) (int, error) {
	s, err := GetEnvStr(key)
	if err != nil {
		return 0, err
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return v, nil
}

func GetEnvInt64(key string) (int64, error) {
	s, err := GetEnvStr(key)
	if err != nil {
		return 0, err
	}
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return v, nil
}

func GetEnvArrayStr(key string) ([]string, error) {
	s, err := GetEnvStr(key)
	if err != nil {
		return nil, err
	}

	items := strings.Split(s, ",")

	for i, item := range items {
		items[i] = strings.Trim(item, " ")
	}
	return items, nil
}

func GetEnvBool(key string) (bool, error) {
	s, err := GetEnvStr(key)
	if err != nil {
		return false, err
	}
	v, err := strconv.ParseBool(s)
	if err != nil {
		return false, err
	}
	return v, nil
}
