// SPDX-License-Identifier: MIT

package ast

import (
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/core/messagetest"
	"github.com/caixw/apidoc/v7/internal/locale"
	"github.com/caixw/apidoc/v7/internal/token"
)

var (
	_ token.Sanitizer = &Param{}
	_ token.Sanitizer = &Request{}
	_ token.Sanitizer = &APIDoc{}
	_ token.Sanitizer = &Path{}
	_ token.Sanitizer = &Enum{}
	_ token.Sanitizer = &XMLNamespace{}
)

func newEmptyParser(a *assert.Assertion) *token.Parser {
	rslt := messagetest.NewMessageHandler()
	p, err := token.NewParser(rslt.Handler, core.Block{})
	a.NotError(err).NotNil(p)
	rslt.Handler.Stop()
	a.Empty(rslt.Errors)
	return p
}

func TestCheckXML(t *testing.T) {
	a := assert.New(t)

	xml := &XML{XMLAttr: &BoolAttribute{Value: Bool{Value: true}}}
	a.Error(checkXML(true, true, xml, newEmptyParser(a)))

	xml = &XML{
		XMLAttr:    &BoolAttribute{Value: Bool{Value: true}},
		XMLWrapped: &Attribute{Value: token.String{Value: "wrapped"}},
	}
	a.Error(checkXML(false, false, xml, newEmptyParser(a)))

	xml = &XML{
		XMLAttr:    &BoolAttribute{Value: Bool{Value: true}},
		XMLExtract: &BoolAttribute{Value: Bool{Value: true}},
	}
	a.Error(checkXML(false, false, xml, newEmptyParser(a)))

	xml = &XML{
		XMLAttr:  &BoolAttribute{Value: Bool{Value: true}},
		XMLCData: &BoolAttribute{Value: Bool{Value: true}},
	}
	a.Error(checkXML(false, false, xml, newEmptyParser(a)))

	xml = &XML{
		XMLWrapped: &Attribute{Value: token.String{Value: "wrapped"}},
	}
	a.Error(checkXML(false, false, xml, newEmptyParser(a)))

	xml = &XML{
		XMLExtract:  &BoolAttribute{Value: Bool{Value: true}},
		XMLNSPrefix: &Attribute{Value: token.String{Value: "p1"}},
	}
	a.Error(checkXML(false, false, xml, newEmptyParser(a)))
}

func TestAPIDoc_checkXMLNamespaces(t *testing.T) {
	a := assert.New(t)

	doc := &APIDoc{
		XMLNamespaces: []*XMLNamespace{
			{
				URN: &Attribute{Value: token.String{Value: "urn1"}},
			},
			{
				URN: &Attribute{Value: token.String{Value: "urn1"}},
			},
		},
	}
	a.ErrorString(doc.checkXMLNamespaces(newEmptyParser(a)), locale.Sprintf(locale.ErrDuplicateValue))

	doc = &APIDoc{
		XMLNamespaces: []*XMLNamespace{
			{
				URN:    &Attribute{Value: token.String{Value: "urn1"}},
				Prefix: &Attribute{Value: token.String{Value: "p1"}},
			},
			{
				URN:    &Attribute{Value: token.String{Value: "urn2"}},
				Prefix: &Attribute{Value: token.String{Value: "p1"}},
			},
		},
	}
	a.ErrorString(doc.checkXMLNamespaces(newEmptyParser(a)), locale.Sprintf(locale.ErrDuplicateValue))

	// 两个 prefix 都为空，也是返回重复的值错误
	doc = &APIDoc{
		XMLNamespaces: []*XMLNamespace{
			{
				URN: &Attribute{Value: token.String{Value: "urn1"}},
			},
			{
				URN: &Attribute{Value: token.String{Value: "urn2"}},
			},
		},
	}
	a.ErrorString(doc.checkXMLNamespaces(newEmptyParser(a)), locale.Sprintf(locale.ErrDuplicateValue))

	// 正常
	doc = &APIDoc{
		XMLNamespaces: []*XMLNamespace{
			{
				URN:    &Attribute{Value: token.String{Value: "urn1"}},
				Prefix: &Attribute{Value: token.String{Value: "p1"}},
			},
			{
				URN:    &Attribute{Value: token.String{Value: "urn2"}},
				Prefix: &Attribute{Value: token.String{Value: "p2"}},
			},
			{
				URN:    &Attribute{Value: token.String{Value: "urn3"}},
				Prefix: &Attribute{Value: token.String{Value: "p3"}},
			},
		},
	}
	a.NotError(doc.checkXMLNamespaces(newEmptyParser(a)))
}

func TestAPI_Sanitize(t *testing.T) {
	a := assert.New(t)

	api := &API{}
	a.NotError(api.Sanitize(newEmptyParser(a)))

	// headers

	api.Headers = []*Param{
		{Type: &TypeAttribute{Value: token.String{Value: TypeString}}},
	}
	a.NotError(api.Sanitize(newEmptyParser(a)))

	api.Headers = append(api.Headers, &Param{
		Type: &TypeAttribute{Value: token.String{Value: TypeObject}},
	})
	a.Error(api.Sanitize(newEmptyParser(a)))

	// servers

	api = &API{
		Servers: []*Element{},
	}
	a.NotError(api.Sanitize(newEmptyParser(a)))

	api.Servers = append(api.Servers, &Element{Content: Content{Value: "s1"}})
	a.NotError(api.Sanitize(newEmptyParser(a)))
	api.Servers = append(api.Servers, &Element{Content: Content{Value: "s1"}})
	a.Error(api.Sanitize(newEmptyParser(a)))

	// tags

	api = &API{
		Tags: []*Element{},
	}
	a.NotError(api.Sanitize(newEmptyParser(a)))

	api.Tags = append(api.Tags, &Element{Content: Content{Value: "s1"}})
	a.NotError(api.Sanitize(newEmptyParser(a)))
	api.Tags = append(api.Tags, &Element{Content: Content{Value: "s1"}})
	a.Error(api.Sanitize(newEmptyParser(a)))
}

func TestXMLnamespace_Sanitize(t *testing.T) {
	a := assert.New(t)

	ns := &XMLNamespace{}
	a.Error(ns.Sanitize(newEmptyParser(a)))

	ns.URN = &Attribute{Value: token.String{Value: "urn"}}
	a.NotError(ns.Sanitize(newEmptyParser(a)))
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

	for i, item := range data {
		err := chkEnumsType(item.t, item.enums, newEmptyParser(a))
		if item.err {
			a.Error(err, "not error at %d", i)
		} else {
			a.NotError(err, "err %s at %d", err, i)
		}
	}
}

func TestAPI_sanitizeTags(t *testing.T) {
	a := assert.New(t)

	api := &API{}
	a.Panic(func() {
		api.sanitizeTags()
	})

	doc := &APIDoc{}
	api.doc = doc
	a.NotError(api.sanitizeTags())

	api.Servers = []*Element{
		{Content: Content{Value: "s1"}},
	}
	a.Error(api.sanitizeTags())
	doc.Servers = []*Server{
		{Name: &Attribute{Value: token.String{Value: "s1"}}},
	}
	a.NotError(api.sanitizeTags())

	api.Tags = []*Element{
		{Content: Content{Value: "t1"}},
	}
	a.Error(api.sanitizeTags())
	doc.Tags = []*Tag{
		{Name: &Attribute{Value: token.String{Value: "t1"}}},
	}
	a.NotError(api.sanitizeTags())
}
