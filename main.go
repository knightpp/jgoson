package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"go/format"
	"io"
	"log"
	"os"
	"reflect"
)

//go:embed testdata.json
var testData []byte

func main() {
	var v any
	err := json.Unmarshal(testData, &v)
	if err != nil {
		log.Fatal(err)
	}

	t := recursion(v)

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(t)

	buf := bytes.Buffer{}
	t.ToGoInline(&buf)

	src, err := format.Source(buf.Bytes())
	if err != nil {
		fmt.Println(buf.String())
		log.Fatal(err)
	}

	fmt.Println(string(src))

	// buf := bytes.Buffer{}
	// printStruct(&buf, Type{
	// 	Name:  "Generated",
	// 	Types: structs,
	// })

	// src, err := format.Source(buf.Bytes())
	// if err != nil {
	// 	fmt.Println(buf.String())
	// 	log.Fatal(err)
	// }

	// fmt.Println(string(src))
}

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
		fmt.Fprintln(w)
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
}

func recursion(value any) *Type {
	return recursionInner(value, "Generated", 0)
}

func recursionInner(value any, parentName string, depth int) *Type {
	t := &Type{
		Name:   parentName,
		Fields: nil,
	}

	if m, ok := value.(map[string]any); ok {
		for k, v := range m {
			t.Fields = append(t.Fields, Field{
				Name: k,
				Type: recursionInner(v, k, depth+1),
			})
		}
		return t
	} else if s, ok := value.([]any); ok {
		if len(s) > 0 {
			newT := recursionInner(s[0], parentName, depth+1)
			newT.IsSlice = true

			t.Fields = append(t.Fields, Field{
				Name: parentName,
				Type: newT,
			})

			return t
		}
		return t
	} else {
		return &Type{
			Name: reflect.TypeOf(value).String(),
		}
	}
}
