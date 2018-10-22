// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package output

import (
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/internal/options"
)

var _ options.Sanitizer = &Options{}

var (
	_ marshaler = apidocJSONMarshal
	_ marshaler = apidocYAMLMarshal
)

func TestOptions_Sanitize(t *testing.T) {
	a := assert.New(t)
	o := &Options{}
	a.Error(o.Sanitize())

	o.Path = "./testdir/apidoc.json"
	a.NotError(o.Sanitize())
	a.Equal(o.marshal, marshaler(apidocJSONMarshal))

	o.Type = typeApidocYAML
	a.NotError(o.Sanitize())
	a.Equal(o.marshal, marshaler(apidocYAMLMarshal))

	o.Type = "unknown"
	a.Error(o.Sanitize())
}

func TestOptions_contains(t *testing.T) {
	a := assert.New(t)

	o := &Options{}
	a.True(o.contains("not exists"))

	o.Groups = []string{"g1", "g2"}
	a.True(o.contains("g1")).
		True(o.contains("g2")).
		False(o.contains("not exists"))
}
