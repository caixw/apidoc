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

func TestParsePath(t *testing.T) {
	a := assert.New(t)

	data := []*struct {
		path   string
		params map[string]struct{}
		err    bool
	}{
		{},
		{
			path: "/path",
		},

		{
			path:   "/{path}",
			params: map[string]struct{}{"path": {}},
		},
		{
			path: "/{path}/{p2}",
			params: map[string]struct{}{
				"path": {},
				"p2":   {},
			},
		},
		{
			path: "/{{path}",
			err:  true,
		},
		{
			path: "/{path}}",
			err:  true,
		},
		{
			path: "/{p{ath}",
			err:  true,
		},
		{
			path: "/{path",
			err:  true,
		},
	}

	for _, item := range data {
		p, err := parsePath(item.path)

		if item.err {
			a.Error(err).Nil(p)
			continue
		}
		a.NotError(err).Equal(p, item.params)
	}
}

func TestChkEnumsType(t *testing.T) {
	a := assert.New(t)

	data := []*struct {
		t     *TypeAttribute
		enums []*Enum
		err   bool
	}{
		{},
		{
			t: &TypeAttribute{Value: token.String{Value: TypeString}},
			enums: []*Enum{
				{Value: &Attribute{Value: token.String{Value: "123"}}},
				{Value: &Attribute{Value: token.String{Value: "str"}}},
			},
		},
		{
			t: &TypeAttribute{Value: token.String{Value: TypeNumber}},
			enums: []*Enum{
				{Value: &Attribute{Value: token.String{Value: "123"}}},
				{Value: &Attribute{Value: token.String{Value: "-1.9"}}},
			},
		},
		{
			t: &TypeAttribute{Value: token.String{Value: TypeBool}},
			enums: []*Enum{
				{Value: &Attribute{Value: token.String{Value: "true"}}},
				{Value: &Attribute{Value: token.String{Value: "0"}}},
			},
		},

		{
			t: &TypeAttribute{Value: token.String{Value: TypeBool}},
			enums: []*Enum{
				{Value: &Attribute{Value: token.String{Value: "string"}}},
				{Value: &Attribute{Value: token.String{Value: "0"}}},
			},
			err: true,
		},
		{
			t: &TypeAttribute{Value: token.String{Value: TypeNumber}},
			enums: []*Enum{
				{Value: &Attribute{Value: token.String{Value: "string"}}},
				{Value: &Attribute{Value: token.String{Value: "0"}}},
			},
			err: true,
		},
		{ // object 是不允许的
			t: &TypeAttribute{Value: token.String{Value: TypeObject}},
			enums: []*Enum{
				{Value: &Attribute{Value: token.String{Value: "string"}}},
				{Value: &Attribute{Value: token.String{Value: "0"}}},
			},
			err: true,
		},
	}

	p, err := token.NewParser(core.Block{})
	a.NotError(err).NotNil(p)

	for i, item := range data {
		err := chkEnumsType(item.t, item.enums, p)
		if item.err {
			a.Error(err, "not error at %d", i)
		} else {
			a.NotError(err, "err %s at %d", err, i)
		}
	}
}
