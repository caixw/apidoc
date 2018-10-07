// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package docs

import (
	"testing"

	"github.com/issue9/assert"
)

func TestBuildSchema(t *testing.T) {
	a := assert.New(t)

	schema := &Schema{}
	a.NotError(buildSchema(schema, nil, []byte("object"), []byte("required"), []byte("desc")))
	a.Equal(schema.Type, "object")

	schema = &Schema{}
	a.NotError(buildSchema(schema, []byte("array"), []byte("array.object"), []byte("required"), []byte("desc")))
	arr := schema.Properties["array"]
	a.NotNil(arr)
	a.Equal(arr.Type, "array")
	a.Equal(arr.Items.Type, "object")
	a.Equal(len(schema.Required), 1).
		Equal(schema.Required[0], "array")

	schema = &Schema{}
	a.NotError(buildSchema(schema, []byte("obj.array"), []byte("array.object"), []byte("required"), []byte("desc")))
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
