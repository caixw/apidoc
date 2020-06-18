// SPDX-License-Identifier: MIT

package openapi

import (
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v7/internal/ast"
	"github.com/caixw/apidoc/v7/internal/token"
)

func TestNewSchema(t *testing.T) {
	a := assert.New(t)

	d := &ast.APIDoc{}
	input := &ast.Param{
		Name:       &ast.Attribute{Value: token.String{Value: "name"}},
		Type:       &ast.TypeAttribute{Value: token.String{Value: ast.TypeBool}},
		Deprecated: &ast.VersionAttribute{Value: token.String{Value: "v1.1.0"}},
		Default:    &ast.Attribute{Value: token.String{Value: "true"}},
		Optional:   &ast.BoolAttribute{Value: ast.Bool{Value: true}},
		Array:      &ast.BoolAttribute{Value: ast.Bool{Value: false}},
		Summary:    &ast.Attribute{Value: token.String{Value: "summary"}},
	}
	output := newSchema(d, input, true)
	a.Equal(output.Type, TypeBool).
		True(output.Deprecated).
		Equal(output.Title, input.Summary.V())

	input.Array = &ast.BoolAttribute{Value: ast.Bool{Value: true}}
	output = newSchema(d, input, true)
	a.Equal(output.Type, TypeArray).
		Equal(output.Items.Type, TypeBool).
		True(output.Items.Deprecated).
		Equal(output.Items.Title, input.Summary.V())

	input.Enums = []*ast.Enum{
		{
			Value:   &ast.Attribute{Value: token.String{Value: "v1"}},
			Summary: &ast.Attribute{Value: token.String{Value: "s1"}},
		},
		{
			Value: &ast.Attribute{Value: token.String{Value: "v2"}},
			Description: &ast.Richtext{
				Text: &ast.CData{Value: token.String{Value: "s2"}},
			},
			Deprecated: &ast.VersionAttribute{Value: token.String{Value: "1.0.1"}},
		},
	}
	output = newSchema(d, input, false)
	a.Equal(output.Type, TypeBool).
		Equal(2, len(output.Enum)).
		Equal(output.Enum, []string{"v1", "v2"})

	input = &ast.Param{
		Type: &ast.TypeAttribute{Value: token.String{Value: ast.TypeNumber}},
		Items: []*ast.Param{
			{
				Name: &ast.Attribute{Value: token.String{Value: "p1"}},
				Type: &ast.TypeAttribute{Value: token.String{Value: ast.TypeString}},
			},
			{
				Name: &ast.Attribute{Value: token.String{Value: "p2"}},
				Type: &ast.TypeAttribute{Value: token.String{Value: ast.TypeNumber}},
			},
		},
	}
	output = newSchema(d, input, true)
	a.Equal(output.Type, "").
		Equal(len(input.Items), len(output.Properties)).
		Equal(output.Properties["p1"].Type, TypeString).
		Equal(output.Properties["p2"].Type, TypeDouble)

	a.NotError(output.sanitize())
}
