// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package schema

import (
	"testing"

	"github.com/issue9/assert"
)

func TestIsOptional(t *testing.T) {
	a := assert.New(t)

	a.False(isOptional(requiredBytes))
	a.True(isOptional([]byte("optional")))
	a.True(isOptional([]byte("")))
}

func TestParseEnum(t *testing.T) {
	a := assert.New(t)

	enums := parseEnum([]byte(`xx
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
