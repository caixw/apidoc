// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package doc

import (
	"testing"

	"github.com/issue9/assert"
)

func TestSchema_build(t *testing.T) {
	a := assert.New(t)

	schema := &Schema{}
	a.NotError(schema.build(nil, []byte("object"), requiredBytes, []byte("desc")))
	a.Equal(schema.Type, "object")

	schema = &Schema{}
	a.NotError(schema.build([]byte("array"), []byte("array.object"), requiredBytes, []byte("desc")))
	arr := schema.Properties["array"]
	a.NotNil(arr)
	a.Equal(arr.Type, Array)
	a.Equal(arr.Items.Type, Object)
	a.Equal(len(schema.Required), 1).
		Equal(schema.Required[0], "array")

	schema = &Schema{}
	a.NotError(schema.build([]byte("obj.array"), []byte("array.object"), requiredBytes, []byte("desc")))
	obj := schema.Properties["obj"]
	a.NotNil(obj)
	arr = obj.Properties["array"]
	a.NotNil(arr)
	a.Equal(arr.Type, "array")
	a.Equal(arr.Items.Type, "object")
	a.Equal(len(obj.Required), 1).
		Equal(obj.Required[0], "array")

	// 可选的参数
	schema = &Schema{}
	a.NotError(schema.build([]byte("array"), []byte("array.object"), []byte("optional"), []byte("desc")))
	arr = schema.Properties["array"]
	a.NotNil(arr)
	a.Equal(arr.Type, Array)
	a.Equal(arr.Items.Type, Object)
	a.Equal(len(schema.Required), 0)
	a.Empty(arr.Default)

	schema = &Schema{}
	a.NotError(schema.build([]byte("string"), []byte("string"), []byte("optional"), []byte("desc")))
	str := schema.Properties["string"]
	a.NotNil(str)
	a.Equal(str.Type, String)
	a.Equal(len(schema.Required), 0).
		Equal(str.Default, "")
}

func TestParseOptional(t *testing.T) {
	a := assert.New(t)

	opt, def, err := parseOptional(String, "", []byte("optional.1"))
	a.NotError(err).
		Equal(def, "1").
		True(opt)

	// 默认值为 ""
	opt, def, err = parseOptional(String, "", []byte("optional"))
	a.NotError(err).
		Equal(def, "").
		True(opt)

	// 默认值为 0
	opt, def, err = parseOptional(Number, "", []byte("optional"))
	a.NotError(err).
		Equal(def, 0).
		True(opt)

	// 默认值为空数组
	opt, def, err = parseOptional(Array, String, []byte("optional"))
	a.NotError(err).
		NotNil(def).
		Empty(def).
		True(opt)

	opt, def, err = parseOptional(Number, "", []byte("optional.1"))
	a.NotError(err).
		Equal(def, 1).
		True(opt)

	opt, def, err = parseOptional(Bool, "", []byte("optional.true"))
	a.NotError(err).
		Equal(def, true).
		True(opt)

	opt, def, err = parseOptional(Array, Number, []byte("optional.[1,2]"))
	a.NotError(err).
		Equal(def, []int{1, 2}).
		True(opt)

	opt, def, err = parseOptional(Array, String, []byte("optional.[1,2]"))
	a.NotError(err).
		Equal(def, []string{"1", "2"}).
		True(opt)

	opt, def, err = parseOptional(Array, String, []byte("required.[1,2]"))
	a.NotError(err).
		Equal(def, []string{"1", "2"}).
		False(opt)
}

func TestParseArray(t *testing.T) {
	a := assert.New(t)

	vals := parseArray([]byte("[a1,a2,a3]"))
	a.Equal(vals, [][]byte{[]byte("a1"), []byte("a2"), []byte("a3")})

	vals = parseArray([]byte("a1,a2,a3"))
	a.Equal(vals, [][]byte{[]byte("a1"), []byte("a2"), []byte("a3")})

	vals = parseArray([]byte("[a1,a2,a3,]"))
	a.Equal(vals, [][]byte{[]byte("a1"), []byte("a2"), []byte("a3"), []byte("")})

	vals = parseArray([]byte("[a1,,a2,a3,  ,]"))
	a.Equal(vals, [][]byte{[]byte("a1"), []byte(""), []byte("a2"), []byte("a3"), []byte(""), []byte("")})

	vals = parseArray([]byte("[a1,a2,  a3  ]"))
	a.Equal(vals, [][]byte{[]byte("a1"), []byte("a2"), []byte("a3")})

	vals = parseArray([]byte("[a1]"))
	a.Equal(vals, [][]byte{[]byte("a1")})

	vals = parseArray([]byte("a1"))
	a.Equal(vals, [][]byte{[]byte("a1")})
}

func TestParseEnum(t *testing.T) {
	a := assert.New(t)

	enums, err := parseEnum(String, []byte(`xx
	- state1 状态 1 描述
	- s2 状态 2 描述
	- s3 状态 3 描述
	状态3 换行描述
	- 状态4 状态 4 描述`))
	a.NotError(err).
		Equal(enums, []string{"state1", "s2", "s3", "状态4"})
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
