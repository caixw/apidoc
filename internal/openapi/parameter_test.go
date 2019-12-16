// SPDX-License-Identifier: MIT

package openapi

import (
	"testing"

	"github.com/issue9/assert"
)

func TestParameter_sanitize(t *testing.T) {
	a := assert.New(t)

	p := &Parameter{}
	a.Error(p.sanitize())

	p.Style = Style{Style: StyleDeepObject}
	a.Error(p.sanitize())

	p.IN = ParameterINPath
	a.NotError(p.sanitize())
}

func TestHeader_sanitize(t *testing.T) {
	a := assert.New(t)

	p := &Header{}
	a.Error(p.sanitize())

	p.Style = Style{Style: StyleDeepObject}
	a.NotError(p.sanitize())

	// IN 只能为空
	p.IN = ParameterINPath
	a.Error(p.sanitize())

	p.IN = ""
	p.Name = "test"
	a.Error(p.sanitize())

	p.Name = ""
	a.NotError(p.sanitize())
}
