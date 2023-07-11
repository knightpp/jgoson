package jgoson

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

func Parse(value map[string]any) *Type {
	return parseRecursive(value, "Generated", 0)
}

func ParseJSON(r io.Reader) (*Type, error) {
	var m map[string]any

	err := json.NewDecoder(r).Decode(&m)
	if err != nil {
		return nil, err
	}

	return Parse(m), nil
}

// Type represents a type in Go.
type Type struct {
	Name    string  `json:"name,omitempty"`
	Fields  []Field `json:"fields,omitempty"`
	IsSlice bool    `json:"is_slice,omitempty"`
}

// Field represents a field in a struct.
type Field struct {
	Name string `json:"name,omitempty"`
	Type *Type  `json:"type,omitempty"`
}

func (t *Type) IsStruct() bool {
	return len(t.Fields) != 0
}

// ToGoInline generates Go code for the given type using the provided config, and the types are inlined.
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

// ToGo generates Go code for the given type using the provided config, but the types are not inlined.
func (t *Type) ToGo(w io.Writer, cfg Config) {
	cfg.FillDefaults()

	t.toGo(w, &cfg)
}

func (t *Type) toGo(w io.Writer, cfg *Config) {
	if !t.IsStruct() {
		fmt.Fprint(w, t.Name)
		return
	}

	types := []*Type{t}
	for i := 0; i < len(types); i++ {
		fmt.Fprintln(w, "type", cfg.StructNameFn(types[i].Name), "struct {")
		for _, field := range types[i].Fields {
			types = append(types, field.toGo(w, cfg)...)
			tagValue := strings.Join(append([]string{cfg.TagNameFn(field.Name)}, cfg.TagOpts...), ",")
			fmt.Fprintf(w, " `%s:\"%s\"`", cfg.Tag, tagValue)
			fmt.Fprintln(w)
		}
		fmt.Fprintln(w, "}")
		fmt.Fprintln(w)
	}
}

func (f Field) toGo(w io.Writer, cfg *Config) (types []*Type) {
	fmt.Fprint(w, cfg.StructFieldFn(f.Name), " ")

	if f.Type.IsSlice {
		fmt.Fprint(w, "[]")
	}

	if f.Type.IsStruct() {
		fmt.Fprint(w, cfg.StructNameFn(f.Type.Name))
		types = append(types, f.Type)
	} else {
		f.Type.toGo(w, cfg)
	}

	return types
}
