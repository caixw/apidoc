// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package doc

import (
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/doc/lexer"
)

func TestBuildSchema(t *testing.T) {
	a := assert.New(t)
	tag := &lexer.Tag{}

	schema := &Schema{}
	a.NotError(buildSchema(tag, schema, nil, []byte("object"), []byte("required"), []byte("desc")))
	a.Equal(schema.Type, "object")

	schema = &Schema{}
	a.NotError(buildSchema(tag, schema, []byte("array"), []byte("array.object"), []byte("required"), []byte("desc")))
	arr := schema.Properties["array"]
	a.NotNil(arr)
	a.Equal(arr.Type, "array")
	a.Equal(arr.Items.Type, "object")
	a.Equal(len(schema.Required), 1).
		Equal(schema.Required[0], "array")

	schema = &Schema{}
	a.NotError(buildSchema(tag, schema, []byte("obj.array"), []byte("array.object"), []byte("required"), []byte("desc")))
	obj := schema.Properties["obj"]
	a.NotNil(obj)
	arr = obj.Properties["array"]
	a.NotNil(arr)
	a.Equal(arr.Type, "array")
	a.Equal(arr.Items.Type, "object")
	a.Equal(len(obj.Required), 1).
		Equal(obj.Required[0], "array")

}

func TestIsRequired(t *testing.T) {
	a := assert.New(t)

	a.True(isRequired("required"))
	a.False(isRequired("optional"))
	a.False(isRequired(""))
}

func TestEnum(t *testing.T) {
	a := assert.New(t)

	enums := enum([]byte(`xx
	- state1 状态 1 描述
	- s2 状态 2 描述
	* s3 状态 3 描述
	状态3 换行描述
	- 状态4 状态 4 描述`))
	a.Equal(enums, []string{"state1", "s2", "s3", "状态4"})
}

func TestConvertEnumType(t *testing.T) {
	a := assert.New(t)

	enums := []string{"1", "2", "3"}
	vals, err := convertEnumType(enums, Integer)
	a.NotError(err).
		Equal(vals, []interface{}{1, 2, 3})

	enums = []string{"1", "2", "3"}
	vals, err = convertEnumType(enums, Number)
	a.NotError(err).
		Equal(vals, []interface{}{1, 2, 3})

	enums = []string{"true", "false", "true"}
	vals, err = convertEnumType(enums, Bool)
	a.NotError(err).
		Equal(vals, []interface{}{true, false, true})
}
