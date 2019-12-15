// SPDX-License-Identifier: MIT

package openapi

import (
	"testing"

	"github.com/issue9/assert"
)

func TestParameter_Sanitize(t *testing.T) {
	a := assert.New(t)

	p := &Parameter{}
	a.Error(p.Sanitize())

	p.Style = Style{Style: StyleDeepObject}
	a.Error(p.Sanitize())

	p.IN = ParameterINPath
	a.NotError(p.Sanitize())
}

func TestHeader_Sanitize(t *testing.T) {
	a := assert.New(t)

	p := &Header{}
	a.Error(p.Sanitize())

	p.Style = Style{Style: StyleDeepObject}
	a.NotError(p.Sanitize())

	// IN 只能为空
	p.IN = ParameterINPath
	a.Error(p.Sanitize())

	p.IN = ""
	p.Name = "test"
	a.Error(p.Sanitize())

	p.Name = ""
	a.NotError(p.Sanitize())
}
