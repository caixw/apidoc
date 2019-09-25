// SPDX-License-Identifier: MIT

package output

import (
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v5/internal/openapi"
)

var (
	_ marshaler = xmlMarshal

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

func TestBuildOptions(t *testing.T) {
	a := assert.New(t)
	oo := &Options{}
	a.Error(oo.sanitize())

	oo.Path = "./testdir/apidoc.json"
	a.NotError(oo.sanitize())
	a.Equal(oo.marshal, marshaler(xmlMarshal))

	a.NotError(oo.sanitize())

	oo.Type = "unknown"
	a.Error(oo.sanitize())
}
