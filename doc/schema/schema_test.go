// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package schema

import (
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/doc/lexer"
)

func TestBuild(t *testing.T) {
	a := assert.New(t)
	tag := &lexer.Tag{}

	schema := &Schema{}
	a.NotError(schema.Build(tag, nil, []byte("object"), requiredBytes, []byte("desc")))
	a.Equal(schema.Type, "object")

	schema = &Schema{}
	a.NotError(schema.Build(tag, []byte("array"), []byte("array.object"), []byte("required"), []byte("desc")))
	arr := schema.Properties["array"]
	a.NotNil(arr)
	a.Equal(arr.Type, Array)
	a.Equal(arr.Items.Type, Object)
	a.Equal(len(schema.Required), 1).
		Equal(schema.Required[0], "array")

	schema = &Schema{}
	a.NotError(schema.Build(tag, []byte("obj.array"), []byte("array.object"), requiredBytes, []byte("desc")))
	obj := schema.Properties["obj"]
	a.NotNil(obj)
	arr = obj.Properties["array"]
	a.NotNil(arr)
	a.Equal(arr.Type, "array")
	a.Equal(arr.Items.Type, "object")
	a.Equal(len(obj.Required), 1).
		Equal(obj.Required[0], "array")

}
