// SPDX-License-Identifier: MIT

package openapi

import (
	"testing"

	"github.com/caixw/apidoc/v6/doc"
	"github.com/issue9/assert"
)

func TestNewSchema(t *testing.T) {
	a := assert.New(t)

	input := &doc.Param{
		Name:       "name",
		Type:       doc.Bool,
		Deprecated: "v1.1.0",
		Default:    "true",
		Optional:   true,
		Array:      false,
		Summary:    "summary",
	}
	output := newSchema(input, true)
	a.Equal(output.Type, TypeBool).
		True(output.Deprecated).
		Equal(output.Title, input.Summary)

	input.Array = true
	output = newSchema(input, true)
	a.Equal(output.Type, TypeArray).
		Equal(output.Items.Type, TypeBool).
		True(output.Items.Deprecated).
		Equal(output.Items.Title, input.Summary)

	input.Enums = []*doc.Enum{
		{
			Value:   "v1",
			Summary: "s1",
		},
		{
			Value:       "v2",
			Description: doc.Richtext{Text: "s2"},
			Deprecated:  "1.0.1",
		},
	}
	output = newSchema(input, false)
	a.Equal(output.Type, TypeBool).
		Equal(output.Enum, []string{"v1", "v2"})

	input = &doc.Param{
		Type: doc.Number,
		Items: []*doc.Param{
			{
				Name: "p1",
				Type: doc.String,
			},
			{
				Name: "p2",
				Type: doc.Number,
			},
		},
	}
	output = newSchema(input, true)
	a.Equal(output.Type, "").
		Equal(len(input.Items), len(output.Properties)).
		Equal(output.Properties["p1"].Type, TypeString).
		Equal(output.Properties["p2"].Type, TypeInt)

	a.NotError(output.sanitize())
}
