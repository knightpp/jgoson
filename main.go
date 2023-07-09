package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"go/format"
	"log"
	"oleksii-brainfuck/jgoson"
	"os"
)

//go:embed testdata.json
var testData []byte

func main() {
	var v map[string]any
	err := json.Unmarshal(testData, &v)
	if err != nil {
		log.Fatal(err)
	}

	t := jgoson.Parse(v)

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(t)

	buf := bytes.Buffer{}
	t.ToGoInline(&buf, jgoson.Config{})

	src, err := format.Source(buf.Bytes())
	if err != nil {
		fmt.Println(buf.String())
		log.Fatal(err)
	}

	fmt.Println(string(src))
}
