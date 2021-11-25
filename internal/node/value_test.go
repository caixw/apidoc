// SPDX-License-Identifier: MIT

package node

import (
	"reflect"
	"testing"

	"github.com/issue9/assert/v2"
)

func TestParseValue(t *testing.T) {
	a := assert.New(t, false)

	v := ParseValue(reflect.ValueOf(intTag{}))
	a.Equal(v.Name, "number").
		Equal(v.Usage, "usage-number").
		False(v.Omitempty)

	v = ParseValue(reflect.ValueOf(struct{}{}))
	a.Nil(v)

	v = ParseValue(reflect.ValueOf(&struct {
		Value int `apidoc:"-"`
	}{}))
	a.Nil(v)

	a.Panic(func() {
		v = ParseValue(reflect.ValueOf(1))
	})
}
