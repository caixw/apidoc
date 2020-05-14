// SPDX-License-Identifier: MIT

package ast

import (
	"testing"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/token"
	"github.com/issue9/assert"
)

var (
	_ token.Sanitizer = &Param{}
	_ token.Sanitizer = &Request{}
	_ token.Sanitizer = &APIDoc{}
	_ token.Sanitizer = &Path{}
	_ token.Sanitizer = &Enum{}
)

func TestAPI_Sanitize(t *testing.T) {
	a := assert.New(t)

	p, err := token.NewParser(core.Block{})
	a.NotError(err).NotNil(p)

	api := &API{}
	a.NotError(api.Sanitize(p))

	api.Headers = []*Param{
		{
			Type: &TypeAttribute{Value: token.String{Value: TypeString}},
		},
	}
	a.NotError(api.Sanitize(p))

	api.Headers = append(api.Headers, &Param{
		Type: &TypeAttribute{Value: token.String{Value: TypeObject}},
	})
	a.Error(api.Sanitize(p))
}
