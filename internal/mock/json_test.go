// SPDX-License-Identifier: MIT

package mock

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v7/internal/ast"
	"github.com/caixw/apidoc/v7/internal/xmlenc"
)

func TestJSONValidator_valid(t *testing.T) {
	a := assert.New(t)

	r := &ast.Request{Type: &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeString}}}
	v := newJSONValidator(r)
	d := json.NewDecoder(strings.NewReader(`"str"`))
	a.NotError(v.valid(d))
	a.Empty(v.names)

	r = &ast.Request{Type: &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeNumber}}}
	v = newJSONValidator(r)
	d = json.NewDecoder(strings.NewReader(`5.0`))
	a.NotError(v.valid(d))
	a.Empty(v.names)

	r = &ast.Request{Type: &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeString}}}
	v = newJSONValidator(r)
	d = json.NewDecoder(strings.NewReader(`5.0`))
	a.Error(v.valid(d))
}

func TestJSONValidator_find(t *testing.T) {
	a := assert.New(t)

	v := &jsonValidator{
		param: (&ast.Request{
			Name: &ast.Attribute{Value: xmlenc.String{Value: "root"}},
			Type: &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeObject}},
			Items: []*ast.Param{
				{
					Type: &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeString}},
					Name: &ast.Attribute{Value: xmlenc.String{Value: "name"}},
				},
				{
					Type: &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeNumber}},
					Name: &ast.Attribute{Value: xmlenc.String{Value: "id"}},
				},
				{
					Type: &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeObject}},
					Name: &ast.Attribute{Value: xmlenc.String{Value: "group"}},
					Items: []*ast.Param{
						{
							Type: &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeString}},
							Name: &ast.Attribute{Value: xmlenc.String{Value: "name"}},
						},
						{
							Type: &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeNumber}},
							Name: &ast.Attribute{Value: xmlenc.String{Value: "id"}},
						},
						{
							Name: &ast.Attribute{Value: xmlenc.String{Value: "tags"}},
							Type: &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeObject}},
							Items: []*ast.Param{
								{
									Type: &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeString}},
									Name: &ast.Attribute{Value: xmlenc.String{Value: "name"}},
								},
								{
									Type: &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeNumber}},
									Name: &ast.Attribute{Value: xmlenc.String{Value: "id"}},
								},
							},
						}, // end tags
					},
				}, // end group
			},
		}).Param(),
	}

	v.names = []string{}
	p := v.find()
	a.Equal(p, v.param)

	v.names = nil
	p = v.find()
	a.Equal(p, v.param)

	v.names = []string{""}
	p = v.find()
	a.Nil(p)

	v.names = []string{"name"}
	p = v.find()
	a.NotNil(p).Equal(p.Type.V(), ast.TypeString)

	v.names = []string{"not-exists"}
	p = v.find()
	a.Nil(p)

	v.names = []string{"group", "id"}
	p = v.find()
	a.NotNil(p).Equal(p.Type.V(), ast.TypeNumber)

	v.names = []string{"group", "tags", "id"}
	p = v.find()
	a.NotNil(p).Equal(p.Type.V(), ast.TypeNumber)
}

func TestValidJSON(t *testing.T) {
	a := assert.New(t)

	for _, item := range data {
		err := validJSON(item.Type, []byte(item.JSON))
		a.NotError(err, "测试 %s 时返回错误值 %s", item.Title, err)
	}
}

func TestBuildJSON(t *testing.T) {
	a := assert.New(t)

	for _, item := range data {
		data, err := buildJSON(item.Type, indent, testOptions)

		a.NotError(err, "测试 %s 返回了错误值 %s", item.Title, err).
			Equal(string(data), item.JSON, "测试 %s 失败 v1:%s,v2:%s", item.Title, string(data), item.JSON)
	}
}
