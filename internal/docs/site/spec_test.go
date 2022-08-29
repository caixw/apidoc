// SPDX-License-Identifier: MIT

package site

import (
	"testing"

	"github.com/issue9/assert/v3"

	"github.com/caixw/apidoc/v7/internal/locale"
)

type (
	objectTag struct {
		RootName struct{}  `apidoc:"apidoc,meta,usage-root"`
		ID       intAttr   `apidoc:"id,attr,usage"`
		Name     stringTag `apidoc:"name,elem,usage"`
	}

	stringTag struct {
		Value    string   `apidoc:"-"`
		RootName struct{} `apidoc:"string,meta,usage-string"`
	}

	intAttr struct {
		Value    int      `apidoc:"-"`
		RootName struct{} `apidoc:"number,meta,usage-number"`
	}
)

func TestNewSpec(t *testing.T) {
	a := assert.New(t, false)

	ts := &doc{}
	err := ts.newSpec(&objectTag{})
	a.NotError(err)
	ts2 := &doc{Spec: []*spec{
		{
			Name:  "apidoc",
			Usage: innerXML{Text: locale.Sprintf("usage-root")},
			Items: []*item{
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
			Usage: innerXML{Text: locale.Sprintf("usage-number")},
			Items: []*item{},
		},
		{
			Name:  "string",
			Usage: innerXML{Text: locale.Sprintf("usage-string")},
			Items: []*item{},
		},
	}}
	a.Equal(ts.Spec, ts2.Spec, "not equal\nv1=%#v\nv2=%#v\n", ts.Spec, ts2.Spec)
}
