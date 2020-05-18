// SPDX-License-Identifier: MIT

// +build ignore

package main

import (
	"bytes"
	"fmt"
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
	return buildConfigObject(doc, "", reflect.TypeOf(build.Config{}))
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

	typeName := t.Kind().String()
	if t.Kind() == reflect.Struct {
		typeName = "object"
	}

	doc.Config = append(doc.Config, &localedoc.Item{
		Name:     name,
		Type:     typeName,
		Array:    array,
		Required: !omitempty,
		Usage:    locale.Sprintf("usage-config-" + name),
	})

	if isPrimitive(t) {
		return nil
	} else if t.Kind() != reflect.Struct {
		panic(fmt.Sprintf("字段 %s 的类型 %s 无法处理", f.Name, t.Kind()))
	}

	return buildConfigObject(doc, name, t)
}

// 调用方需要保证 t.Kind() 为 reflect.Struct
func buildConfigObject(doc *localedoc.LocaleDoc, parent string, t reflect.Type) error {
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if !unicode.IsUpper(rune(f.Name[0])) || f.Tag.Get("yaml") == "-" {
			continue
		}
		if err := buildConfigItem(doc, parent, f); err != nil {
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
