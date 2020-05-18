// SPDX-License-Identifier: MIT

package token

import (
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v7/internal/docs/localedoc"
	"github.com/caixw/apidoc/v7/internal/locale"
)

func TestNewTypes(t *testing.T) {
	a := assert.New(t)

	ts := &localedoc.LocaleDoc{}
	err := NewTypes(ts, &objectTag{})
	a.NotError(err)
	ts2 := &typeList{Types: []*localedoc.Type{
		{
			Name:  "apidoc",
			Usage: localedoc.InnerXML{Text: locale.Sprintf("usage-root")},
			Items: []*localedoc.Item{
				{
					Name:     "@id",
					Usage:    locale.Sprintf("usage"),
					Type:     "number",
					Array:    false,
					Required: true,
				},
				{
					Name:     "name",
					Usage:    locale.Sprintf("usage"),
					Type:     "string",
					Array:    false,
					Required: true,
				},
			},
		},
		{
			Name:  "number",
			Usage: localedoc.InnerXML{Text: locale.Sprintf("usage-number")},
			Items: []*localedoc.Item{},
		},
		{
			Name:  "string",
			Usage: localedoc.InnerXML{Text: locale.Sprintf("usage-string")},
			Items: []*localedoc.Item{},
		},
	}}
	a.Equal(ts.Types, ts2.Types, "not equal\nv1=%#v\nv2=%#v\n", ts.Types, ts2.Types)
}
