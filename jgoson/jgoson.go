package jgoson

import (
	"fmt"
	"io"
	"reflect"
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

func (t *Type) ToGoInline(w io.Writer) {
	fmt.Fprintln(w, "type", t.Name, "struct {")
	for _, field := range t.Fields {
		field.toGoInline(w)
	}
	fmt.Fprintln(w, "}")
}

func (t *Type) toGoInline(w io.Writer) {
	if !t.IsStruct() {
		fmt.Fprint(w, t.Name)
		return
	}

	for _, field := range t.Fields {
		field.toGoInline(w)
	}
}

func (t Field) toGoInline(w io.Writer) {
	fmt.Fprint(w, t.Name, " ")

	if t.Type.IsSlice {
		fmt.Fprint(w, "[]")
	}

	if t.Type.IsStruct() {
		fmt.Fprintln(w, "struct{")
		t.Type.toGoInline(w)
		fmt.Fprintln(w, "}")
		return
	}

	t.Type.toGoInline(w)
	fmt.Fprintln(w)
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
