package jgoson

import (
	"fmt"
	"io"
	"reflect"
	"strings"
)

type Type struct {
	Name    string  `json:"name,omitempty"`
	Fields  []Field `json:"fields,omitempty"`
	IsSlice bool    `json:"is_slice,omitempty"`
}

type Field struct {
	Name string `json:"name,omitempty"`
	Type *Type  `json:"type,omitempty"`
}

func (t *Type) IsStruct() bool {
	return len(t.Fields) != 0
}

func (t *Type) ToGoInline(w io.Writer, cfg Config) {
	cfg.FillDefaults()

	fmt.Fprintln(w, "type", cfg.StructNameFn(t.Name), "struct {")
	t.toGoInline(w, &cfg)
	fmt.Fprintln(w, "}")
}

func (t *Type) toGoInline(w io.Writer, cfg *Config) {
	if !t.IsStruct() {
		fmt.Fprint(w, t.Name)
		return
	}

	for _, field := range t.Fields {
		field.toGoInline(w, cfg)
		tagValue := strings.Join(append([]string{cfg.TagNameFn(field.Name)}, cfg.TagOpts...), ",")
		fmt.Fprintf(w, " `%s:\"%s\"`", cfg.Tag, tagValue)
		fmt.Fprintln(w)
	}
}

func (t Field) toGoInline(w io.Writer, cfg *Config) {
	fmt.Fprint(w, cfg.StructFieldFn(t.Name), " ")

	if t.Type.IsSlice {
		fmt.Fprint(w, "[]")
	}

	if t.Type.IsStruct() {
		fmt.Fprintln(w, "struct{")
		t.Type.toGoInline(w, cfg)
		fmt.Fprint(w, "}")
	} else {
		t.Type.toGoInline(w, cfg)
	}
}

func Parse(value map[string]any) *Type {
	return recursionInner(value, "Generated", 0)
}

func recursionInner(value any, parentName string, depth int) *Type {
	if m, ok := value.(map[string]any); ok {
		t := &Type{
			Name:   parentName,
			Fields: nil,
		}
		for k, v := range m {
			t.Fields = append(t.Fields, Field{
				Name: k,
				Type: recursionInner(v, k, depth+1),
			})
		}
		return t
	} else if s, ok := value.([]any); ok {
		if len(s) == 0 {
			return &Type{
				Name:    "any",
				IsSlice: true,
			}

		}

		newT := recursionInner(s[0], parentName, depth+1)
		newT.IsSlice = true
		return newT
	} else {
		return &Type{
			Name: reflect.TypeOf(value).String(),
		}
	}
}
