package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"log"
	"os"
	"strings"

	"github.com/knightpp/jgoson"
)

var (
	useInline bool
	tag       string
	tagOpts   string
)

func init() {
	flag.BoolVar(&useInline, "inline", true, "inline structs")
	flag.StringVar(&tag, "tag", "json", "tag to use")
	flag.StringVar(&tagOpts, "tag-opts", "omitempty", "tag additional option, comma separated")
}

func main() {
	flag.Parse()

	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	t, err := jgoson.ParseJSON(os.Stdin)
	if err != nil {
		return fmt.Errorf("parse json: %w", err)
	}

	buf := bytes.Buffer{}

	opts := make([]string, 0)
	if tagOpts != "" {
		opts = strings.Split(tagOpts, ",")
	}

	fn := t.ToGo
	if useInline {
		fn = t.ToGoInline
	}
	fn(&buf, jgoson.Config{
		Tag:     tag,
		TagOpts: opts,
	})

	src, err := format.Source(buf.Bytes())
	if err != nil {
		return fmt.Errorf("format source: %w", err)
	}

	fmt.Println(string(src))

	return nil
}
