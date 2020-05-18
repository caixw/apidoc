// SPDX-License-Identifier: MIT

package token

import (
	"testing"

	"github.com/issue9/assert"
	"golang.org/x/text/language"

	"github.com/caixw/apidoc/v7/internal/docs/localedoc"
	"github.com/caixw/apidoc/v7/internal/locale"
)

func TestNewTypes(t *testing.T) {
	a := assert.New(t)
	id := locale.DefaultLocaleID

	ts := &localedoc.LocaleDoc{}
	err := NewTypes(ts, &objectTag{}, language.MustParse(id))
	a.NotError(err)
	ts2 := &typeList{Types: []*localedoc.Type{
		{
			Name:  "apidoc",
			Usage: localedoc.InnerXML{Text: locale.Translate(id, "usage-root")},
			Items: []*localedoc.Item{
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
			Name:  "number",
			Usage: localedoc.InnerXML{Text: locale.Translate(id, "usage-number")},
			Items: []*localedoc.Item{},
		},
		{
			Name:  "string",
			Usage: localedoc.InnerXML{Text: locale.Translate(id, "usage-string")},
			Items: []*localedoc.Item{},
		},
	}}
	a.Equal(ts.Types, ts2.Types, "not equal\nv1=%#v\nv2=%#v\n", ts.Types, ts2.Types)
}
