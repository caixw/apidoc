// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package output

import (
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/internal/options"
	"github.com/caixw/apidoc/output/openapi"
)

var _ options.Sanitizer = &Options{}

var (
	_ marshaler = apidocJSONMarshal
	_ marshaler = apidocYAMLMarshal

	_ marshaler = openapi.JSON
	_ marshaler = openapi.YAML
)

func TestOptions_contains(t *testing.T) {
	a := assert.New(t)

	o := &Options{}
	a.True(o.contains("tag"))
	a.True(o.contains(""))

	o.Tags = []string{"t1", "t2"}
	a.True(o.contains("t1"))
	a.False(o.contains("not-exists"))
	a.False(o.contains(""))
}

func TestOptions_Sanitize(t *testing.T) {
	a := assert.New(t)
	o := &Options{}
	a.Error(o.Sanitize())

	o.Path = "./testdir/apidoc.json"
	a.NotError(o.Sanitize())
	a.Equal(o.marshal, marshaler(apidocJSONMarshal))

	o.Type = ApidocYAML
	a.NotError(o.Sanitize())
	a.Equal(o.marshal, marshaler(apidocYAMLMarshal))

	o.Type = "unknown"
	a.Error(o.Sanitize())
}
