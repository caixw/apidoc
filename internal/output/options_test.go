// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package output

import (
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/internal/output/openapi"
	opt "github.com/caixw/apidoc/options"
)

var (
	_ marshaler = apidocJSONMarshal
	_ marshaler = apidocYAMLMarshal

	_ marshaler = openapi.JSON
	_ marshaler = openapi.YAML
)

func TestOptions_contains(t *testing.T) {
	a := assert.New(t)

	o := &options{}
	a.True(o.contains("tag"))
	a.True(o.contains(""))

	o.Tags = []string{"t1", "t2"}
	a.True(o.contains("t1"))
	a.False(o.contains("not-exists"))
	a.False(o.contains(""))
}

func TestBuildOptions(t *testing.T) {
	a := assert.New(t)
	oo := &opt.Output{}
	o, err := buildOptions(oo)
	a.Error(err).Nil(o)

	oo.Path = "./testdir/apidoc.json"
	o, err = buildOptions(oo)
	a.NotError(err).NotNil(o)
	a.Equal(o.marshal, marshaler(apidocJSONMarshal))

	oo.Type = opt.ApidocYAML
	o, err = buildOptions(oo)
	a.NotError(err).NotNil(o)
	a.Equal(o.marshal, marshaler(apidocYAMLMarshal))

	oo.Type = "unknown"
	o, err = buildOptions(oo)
	a.Error(err).Nil(o)
}
