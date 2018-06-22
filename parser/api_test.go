// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package parser

import (
	"testing"

	"github.com/caixw/apidoc/openapi"
	"github.com/issue9/assert"
)

func TestBuildSchema(t *testing.T) {
	a := assert.New(t)

	schema := &openapi.Schema{}
	a.NotError(buildSchema(schema, nil, []byte("object"), []byte("optional"), []byte("desc")))
	a.Equal(schema.Type, "object")

	schema = &openapi.Schema{}
	a.NotError(buildSchema(schema, []byte("array.list"), []byte("array.object"), []byte("required"), []byte("desc")))
	a.Equal(schema.Properties["array"].Type, "array")
	a.Equal(schema.Properties["array"].Items.Type, "object")
}
