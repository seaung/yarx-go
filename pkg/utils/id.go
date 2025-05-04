package utils

import (
	"strings"

	shortid "github.com/jasonsoft/go-short-id"
)

// 生成唯一ID号
func GenShortUUID() string {
	opts := shortid.Options{
		Number:        4,
		StartWithYear: true,
		EndWithHost:   false,
	}

	return strings.ToLower(shortid.Generate(opts))
}
