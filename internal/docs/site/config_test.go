// SPDX-License-Identifier: MIT

package site

import (
	"reflect"
	"testing"

	"github.com/issue9/assert/v3"
)

func TestParseTag(t *testing.T) {
	a := assert.New(t, false)

	name, omitempty := parseTag(reflect.StructField{Name: "F1"})
	a.Equal(name, "F1").False(omitempty)

	name, omitempty = parseTag(reflect.StructField{Name: "F1", Tag: reflect.StructTag(`apidoc:"xx"`)})
	a.Equal(name, "F1").False(omitempty)

	name, omitempty = parseTag(reflect.StructField{Name: "F1", Tag: reflect.StructTag(`yaml:"xx"`)})
	a.Equal(name, "xx").False(omitempty)

	name, omitempty = parseTag(reflect.StructField{Name: "F1", Tag: reflect.StructTag(`yaml:"xx,omitempty"`)})
	a.Equal(name, "xx").True(omitempty)

	name, omitempty = parseTag(reflect.StructField{Name: "F1", Tag: reflect.StructTag(`yaml:",omitempty"`)})
	a.Equal(name, "F1").True(omitempty)

	name, omitempty = parseTag(reflect.StructField{Name: "F1", Tag: reflect.StructTag(`yaml:"xx,omit"`)})
	a.Equal(name, "xx").False(omitempty)
}
