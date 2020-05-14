// SPDX-License-Identifier: MIT

package token

import (
	"testing"

	"github.com/issue9/assert"
	"golang.org/x/text/language"

	"github.com/caixw/apidoc/v7/internal/locale"
	"github.com/caixw/apidoc/v7/internal/vars"
)

func TestNewTypes(t *testing.T) {
	a := assert.New(t)
	id := vars.DefaultLocaleID

	ts, err := NewTypes(&objectTag{}, language.MustParse(id))
	a.NotError(err).NotNil(ts)
	ts2 := &Types{Types: []*Type{
		{
			Name:  "apidoc",
			Usage: locale.Translate(id, "usage-root"),
			Items: []*Item{
				{
					Name:     "@id",
					Usage:    locale.Translate(id, "usage"),
					Type:     "number",
					Array:    false,
					Required: true,
				},
				{
					Name:     "name",
					Usage:    locale.Translate(id, "usage"),
					Type:     "string",
					Array:    false,
					Required: true,
				},
			},
		},
		{
			Name:  "string",
			Usage: locale.Translate(id, "usage-string"),
			Items: []*Item{},
		},
	}}
	a.Equal(len(ts.Types), len(ts2.Types))
	a.Equal(ts, ts2, "not equal\nv1=%#v\nv2=%#v\n", ts, ts2)
}
