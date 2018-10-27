// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package doc

import (
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/doc/schema"
)

func TestNewParam(t *testing.T) {
	a := assert.New(t)

	p, err := newParam(newTag("name string required  名称"))
	a.NotError(err).
		NotNil(p).
		Equal(p.Name, "name").
		Equal(p.Type.Type, schema.String).
		False(p.Optional).
		Equal(p.Summary, "名称")

	p, err = newParam(newTag("name string optional.v1  名称"))
	a.NotError(err).
		NotNil(p).
		Equal(p.Name, "name").
		Equal(p.Type.Type, schema.String).
		True(p.Optional).
		Equal(p.Summary, "名称")

	p, err = newParam(newTag("name string optional  名称"))
	a.NotError(err).
		NotNil(p).
		Equal(p.Name, "name").
		Equal(p.Type.Type, schema.String).
		True(p.Optional).
		Equal(p.Summary, "名称")
}
