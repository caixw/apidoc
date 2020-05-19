// SPDX-License-Identifier: MIT

package site

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"

	"github.com/caixw/apidoc/v7/build"
	"github.com/caixw/apidoc/v7/internal/locale"
)

func (d *doc) newConfig() error {
	return d.buildConfigObject("", reflect.TypeOf(build.Config{}))
}

func (d *doc) buildConfigItem(parent string, f reflect.StructField) error {
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

	d.Config = append(d.Config, &item{
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

	return d.buildConfigObject(name, t)
}

// 调用方需要保证 t.Kind() 为 reflect.Struct
func (d *doc) buildConfigObject(parent string, t reflect.Type) error {
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if !unicode.IsUpper(rune(f.Name[0])) || f.Tag.Get("yaml") == "-" {
			continue
		}
		if err := d.buildConfigItem(parent, f); err != nil {
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
		return getName(prop[0], f), false
	}

	return getName(prop[0], f), strings.TrimSpace(prop[1]) == "omitempty"
}

func getName(n string, f reflect.StructField) string {
	if n = strings.TrimSpace(n); n != "" {
		return n
	}
	return f.Name
}
