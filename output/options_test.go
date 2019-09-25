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

func TestOptions_sanitize(t *testing.T) {
	a := assert.New(t)

	var o *Options
	a.Error(o.sanitize())

	o = &Options{}
	a.Error(o.sanitize())

	o.Path = "./testdir/apidoc.json"
	a.NotError(o.sanitize())
	a.Equal(o.marshal, marshaler(xmlMarshal))

	o.Type = "unknown"
	a.Error(o.sanitize())
}
