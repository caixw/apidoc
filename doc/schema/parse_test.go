// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package schema

import (
	"testing"

	"github.com/issue9/assert"
)

func TestParseOptional(t *testing.T) {
	a := assert.New(t)

	opt, def, err := parseOptional(String, "", []byte("optional.1"))
	a.NotError(err).
		Equal(def, "1").
		True(opt)

	opt, def, err = parseOptional(String, "", []byte("optional"))
	a.NotError(err).
		Nil(def).
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

func TestIsOptional(t *testing.T) {
	a := assert.New(t)

	a.False(isOptional(requiredBytes))
	a.True(isOptional([]byte("optional")))
	a.True(isOptional([]byte("")))
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
}

func TestParseEnum(t *testing.T) {
	a := assert.New(t)

	enums, err := parseEnum(String, []byte(`xx
	- state1 状态 1 描述
	- s2 状态 2 描述
	* s3 状态 3 描述
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
