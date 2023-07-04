package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"go/format"
	"io"
	"log"
	"reflect"
	"strconv"
)

//go:embed testdata.json
var testData []byte

func main() {
	var v any
	err := json.Unmarshal(testData, &v)
	if err != nil {
		log.Fatal(err)
	}

	structs := recursion(v, nil, 1)

	fmt.Printf("%+#v\n", structs)

	buf := bytes.Buffer{}
	printStruct(&buf, Struct{
		Name:   "Generated",
		Fields: structs,
	})

	src, err := format.Source(buf.Bytes())
	if err != nil {
		fmt.Println(buf.String())
		log.Fatal(err)
	}

	fmt.Println(string(src))
}

func printStruct(w io.Writer, t Struct) {
	fmt.Fprintln(w, "type", t.Name, "struct {")
	printFields(w, t.Fields)
	fmt.Fprintln(w, "}")
}

func printFields(w io.Writer, fields []Struct) {
	for _, field := range fields {
		if len(field.Fields) == 0 {
			fmt.Fprintf(w, "%s %s `%s`\n", field.Name, field.Type, field.Annotation)
			continue
		}

		fmt.Fprintf(w, "%s %s `%s`", field.Name, "GeneratedType", field.Annotation)

		field := field
		printStruct(w, Struct{
			Name:   "GeneratedType",
			Fields: field.Fields,
		})
	}
}

type Struct struct {
	Name       string
	Annotation string
	Fields     []Field
}

type Field struct {
	Name       string
	Type       string
	Annotation string
}

func recursion(value any, in []Struct, depth int) []Struct {
	if m, ok := value.(map[string]any); ok {
		for k, v := range m {
			in = append(in, Struct{
				Name:       k,
				Type:       reflect.TypeOf(v).String(),
				Fields:     recursion(v, nil, depth+1),
				Annotation: `json:"` + k + `"`,
			})
		}
	} else if s, ok := value.([]any); ok {
		for _, v := range s {
			fields := recursion(v, in, depth+1)
			in = append(in, Struct{
				Name:       "[]" + strconv.Itoa(depth),
				Type:       reflect.TypeOf(v).String(),
				Fields:     fields,
				Annotation: "",
			})
		}
	}

	return in
}
