// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package schema

import (
	"testing"

	"github.com/issue9/assert"
)

func TestBuild(t *testing.T) {
	a := assert.New(t)

	schema := &Schema{}
	a.NotError(schema.Build(nil, []byte("object"), requiredBytes, []byte("desc")))
	a.Equal(schema.Type, "object")

	schema = &Schema{}
	a.NotError(schema.Build([]byte("array"), []byte("array.object"), requiredBytes, []byte("desc")))
	arr := schema.Properties["array"]
	a.NotNil(arr)
	a.Equal(arr.Type, Array)
	a.Equal(arr.Items.Type, Object)
	a.Equal(len(schema.Required), 1).
		Equal(schema.Required[0], "array")

	schema = &Schema{}
	a.NotError(schema.Build([]byte("obj.array"), []byte("array.object"), requiredBytes, []byte("desc")))
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
	a.NotError(schema.Build([]byte("array"), []byte("array.object"), []byte("optional"), []byte("desc")))
	arr = schema.Properties["array"]
	a.NotNil(arr)
	a.Equal(arr.Type, Array)
	a.Equal(arr.Items.Type, Object)
	a.Equal(len(schema.Required), 0)
	a.Empty(arr.Default)

	schema = &Schema{}
	a.NotError(schema.Build([]byte("string"), []byte("string"), []byte("optional"), []byte("desc")))
	str := schema.Properties["string"]
	a.NotNil(str)
	a.Equal(str.Type, String)
	a.Equal(len(schema.Required), 0).
		Equal(str.Default, "")
}
