package main

import (
	_ "embed"
	"reflect"
)

//go:embed testdata.json
var testData []byte

func main() {
	// var v any
	// err := json.Unmarshal(testData, &v)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// structs := recursion(v, nil, 1)

	// fmt.Printf("%+#v\n", structs)

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

// func printStruct(w io.Writer, t Type) {
// 	fmt.Fprintln(w, "type", t.Name, "struct {")
// 	printFields(w, t.Types)
// 	fmt.Fprintln(w, "}")
// }

// func printFields(w io.Writer, fields []Type) {
// 	for _, field := range fields {
// 		if len(field.Types) == 0 {
// 			fmt.Fprintf(w, "%s %s `%s`\n", field.Name, field.Type, field.Annotation)
// 			continue
// 		}

// 		fmt.Fprintf(w, "%s %s `%s`", field.Name, "GeneratedType", field.Annotation)

// 		field := field
// 		printStruct(w, Type{
// 			Name:  "GeneratedType",
// 			Types: field.Types,
// 		})
// 	}
// }

type Type struct {
	Name   string
	Fields []Field
}

type Field struct {
	Name string
	Type *Type
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
			t.Fields = append(t.Fields, Field{
				Name: parentName,
				Type: recursionInner(s[0], parentName, depth+1),
			})
		}
		return t
	} else {
		return &Type{
			Name: reflect.TypeOf(value).String(),
		}
	}
}
