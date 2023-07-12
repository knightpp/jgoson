package jgoson

import (
	"strings"
	"unicode"
)

type Config struct {
	Tag       string
	TagNameFn func(fieldName string) string
	TagOpts   []string

	StructNameFn  func(structName string) string
	StructFieldFn func(fieldName string) string
}

func (c *Config) FillDefaults() {
	if c.Tag == "" {
		c.Tag = "json"
	}

	if c.TagNameFn == nil {
		c.TagNameFn = LowerCamelCaseToSnakeCase
	}

	// check for nil to differentiate between empty and nil
	if c.TagOpts == nil {
		c.TagOpts = []string{"omitempty"}
	}

	if c.StructNameFn == nil {
		c.StructNameFn = SnakeCaseToUpperCamelCase
	}

	if c.StructFieldFn == nil {
		c.StructFieldFn = SnakeCaseToUpperCamelCase
	}
}

func LowerCamelCaseToSnakeCase(name string) string {
	var buf strings.Builder
	buf.Grow(len(name))

	for _, c := range name {
		if unicode.IsUpper(c) {
			buf.WriteRune('_')
			buf.WriteRune(unicode.ToLower(c))
		} else {
			buf.WriteRune(c)
		}
	}

	return buf.String()
}

func SnakeCaseToUpperCamelCase(name string) string {
	var buf strings.Builder
	size := 0
	for _, c := range name {
		if c == '_' {
			continue
		}

		size += 4
	}
	buf.Grow(size)

	upperNext := true
	for _, c := range name {
		if c == '_' {
			upperNext = true
			continue
		}

		if upperNext {
			upperNext = false
			c = unicode.ToUpper(c)
		}

		buf.WriteRune(c)
	}

	return buf.String()
}
