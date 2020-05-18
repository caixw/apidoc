// SPDX-License-Identifier: MIT

// +build ignore

package main

import (
	"bytes"
	"io"
	"reflect"
	"strings"
	"unicode"

	"github.com/caixw/apidoc/v7/build"
	"github.com/caixw/apidoc/v7/internal/ast"
	"github.com/caixw/apidoc/v7/internal/cmd"
	"github.com/caixw/apidoc/v7/internal/docs"
	"github.com/caixw/apidoc/v7/internal/docs/localedoc"
	"github.com/caixw/apidoc/v7/internal/docs/makeutil"
	"github.com/caixw/apidoc/v7/internal/locale"
	"github.com/caixw/apidoc/v7/internal/token"
)

func main() {
	for _, tag := range locale.Tags() {
		locale.SetTag(tag)

		doc := &localedoc.LocaleDoc{}
		makeutil.PanicError(makeCommands(doc))
		makeutil.PanicError(makeConfig(doc))
		makeutil.PanicError(token.NewTypes(doc, &ast.APIDoc{}))

		target := docs.Dir().Append(localedoc.Path(tag))
		makeutil.PanicError(makeutil.WriteXML(target, doc, "\t"))
	}
}

func makeCommands(doc *localedoc.LocaleDoc) error {
	out := new(bytes.Buffer)
	opt := cmd.Init(out)
	names := opt.Commands()

	for _, name := range names {
		out.Reset()
		if err := opt.Exec([]string{"help", name}); err != nil {
			return err
		}

		usage, err := out.ReadString('\n')
		if err != nil && err != io.EOF {
			return err
		}

		if usage[len(usage)-1] == '\n' { // 去掉换行符
			usage = usage[:len(usage)-1]
		}
		doc.Commands = append(doc.Commands, &localedoc.Command{
			Name:  name,
			Usage: usage,
		})
	}

	return nil
}

func makeConfig(doc *localedoc.LocaleDoc) error {
	t := reflect.TypeOf(build.Config{})
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if !unicode.IsUpper(rune(f.Name[0])) || f.Tag.Get("yaml") == "-" {
			continue
		}
		if err := buildConfigItem(doc, "", f); err != nil {
			return err
		}
	}
	return nil
}

func buildConfigItem(doc *localedoc.LocaleDoc, parent string, f reflect.StructField) error {
	name, omitempty := parseTag(f)
	if parent != "" {
		name = parent + "." + name
	}

	t := f.Type
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	var array bool
	if t.Kind() == reflect.Array || t.Kind() == reflect.Slice {
		array = true
		t = t.Elem()
		for t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
	}

	if isPrimitive(t) {
		doc.Config = append(doc.Config, &localedoc.Item{
			Name:     name,
			Type:     t.Kind().String(),
			Array:    array,
			Required: !omitempty,
			Usage:    locale.Sprintf("usage-config-" + name),
		})
		return nil
	}

	for i := 0; i < t.NumField(); i++ {
		ff := t.Field(i)
		if !unicode.IsUpper(rune(ff.Name[0])) || ff.Tag.Get("yaml") == "-" {
			continue
		}
		if err := buildConfigItem(doc, name, ff); err != nil {
			return err
		}
	}

	return nil
}

func isPrimitive(t reflect.Type) bool {
	return t.Kind() == reflect.String || (t.Kind() >= reflect.Bool && t.Kind() <= reflect.Complex128)
}

func parseTag(f reflect.StructField) (string, bool) {
	tag := f.Tag.Get("yaml")
	if tag == "" {
		return f.Name, false
	}

	prop := strings.Split(tag, ",")
	if len(prop) == 1 {
		return strings.TrimSpace(prop[0]), false
	}

	return strings.TrimSpace(prop[0]), strings.TrimSpace(prop[1]) == "omitempty"
}
