// SPDX-License-Identifier: MIT

package ast

import (
	"net/http"
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/core/messagetest"
	"github.com/caixw/apidoc/v7/internal/locale"
	"github.com/caixw/apidoc/v7/internal/xmlenc"
)

var (
	_ xmlenc.Sanitizer = &Param{}
	_ xmlenc.Sanitizer = &Request{}
	_ xmlenc.Sanitizer = &APIDoc{}
	_ xmlenc.Sanitizer = &Path{}
	_ xmlenc.Sanitizer = &Enum{}
	_ xmlenc.Sanitizer = &XMLNamespace{}
)

func newEmptyParser(a *assert.Assertion) *xmlenc.Parser {
	rslt := messagetest.NewMessageHandler()
	p, err := xmlenc.NewParser(rslt.Handler, core.Block{})
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
		XMLWrapped: &Attribute{Value: xmlenc.String{Value: "wrapped"}},
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
		XMLWrapped: &Attribute{Value: xmlenc.String{Value: "wrapped"}},
	}
	a.Error(checkXML(false, false, xml, newEmptyParser(a)))

	xml = &XML{
		XMLExtract:  &BoolAttribute{Value: Bool{Value: true}},
		XMLNSPrefix: &Attribute{Value: xmlenc.String{Value: "p1"}},
	}
	a.Error(checkXML(false, false, xml, newEmptyParser(a)))
}

func TestAPIDoc_checkXMLNamespaces(t *testing.T) {
	a := assert.New(t)

	doc := &APIDoc{
		XMLNamespaces: []*XMLNamespace{
			{
				URN: &Attribute{Value: xmlenc.String{Value: "urn1"}},
			},
			{
				URN: &Attribute{Value: xmlenc.String{Value: "urn1"}},
			},
		},
	}
	a.ErrorString(doc.checkXMLNamespaces(newEmptyParser(a)), locale.Sprintf(locale.ErrDuplicateValue))

	doc = &APIDoc{
		XMLNamespaces: []*XMLNamespace{
			{
				URN:    &Attribute{Value: xmlenc.String{Value: "urn1"}},
				Prefix: &Attribute{Value: xmlenc.String{Value: "p1"}},
			},
			{
				URN:    &Attribute{Value: xmlenc.String{Value: "urn2"}},
				Prefix: &Attribute{Value: xmlenc.String{Value: "p1"}},
			},
		},
	}
	a.ErrorString(doc.checkXMLNamespaces(newEmptyParser(a)), locale.Sprintf(locale.ErrDuplicateValue))

	// 两个 prefix 都为空，也是返回重复的值错误
	doc = &APIDoc{
		XMLNamespaces: []*XMLNamespace{
			{
				URN: &Attribute{Value: xmlenc.String{Value: "urn1"}},
			},
			{
				URN: &Attribute{Value: xmlenc.String{Value: "urn2"}},
			},
		},
	}
	a.ErrorString(doc.checkXMLNamespaces(newEmptyParser(a)), locale.Sprintf(locale.ErrDuplicateValue))

	// 正常
	doc = &APIDoc{
		XMLNamespaces: []*XMLNamespace{
			{
				URN:    &Attribute{Value: xmlenc.String{Value: "urn1"}},
				Prefix: &Attribute{Value: xmlenc.String{Value: "p1"}},
			},
			{
				URN:    &Attribute{Value: xmlenc.String{Value: "urn2"}},
				Prefix: &Attribute{Value: xmlenc.String{Value: "p2"}},
			},
			{
				URN:    &Attribute{Value: xmlenc.String{Value: "urn3"}},
				Prefix: &Attribute{Value: xmlenc.String{Value: "p3"}},
			},
		},
	}
	a.NotError(doc.checkXMLNamespaces(newEmptyParser(a)))
}

func TestAPI_Sanitize(t *testing.T) {
	a := assert.New(t)

	api := &API{}
	p, rslt := newParser(a, "", "")
	api.Sanitize(p)
	rslt.Handler.Stop()
	a.Empty(rslt.Errors).Empty(rslt.Warns)

	// headers

	api.Headers = []*Param{
		{Type: &TypeAttribute{Value: xmlenc.String{Value: TypeString}}},
	}
	p, rslt = newParser(a, "", "")
	api.Sanitize(p)
	rslt.Handler.Stop()
	a.Empty(rslt.Errors).Empty(rslt.Warns)

	api.Headers = append(api.Headers, &Param{
		Type: &TypeAttribute{Value: xmlenc.String{Value: TypeObject}},
	})
	p, rslt = newParser(a, "", "")
	api.Sanitize(p)
	rslt.Handler.Stop()
	a.NotEmpty(rslt.Errors)

	// servers

	api = &API{
		Servers: []*ServerValue{},
	}
	p, rslt = newParser(a, "", "")
	api.Sanitize(p)
	rslt.Handler.Stop()
	a.Empty(rslt.Errors)

	api.Servers = append(api.Servers, &ServerValue{Content: Content{Value: "s1"}})
	p, rslt = newParser(a, "", "")
	api.Sanitize(p)
	rslt.Handler.Stop()
	a.Empty(rslt.Errors)

	api.Servers = append(api.Servers, &ServerValue{Content: Content{Value: "s1"}})
	p, rslt = newParser(a, "", "")
	api.Sanitize(p)
	rslt.Handler.Stop()
	a.NotEmpty(rslt.Errors)

	// tags

	api = &API{
		Tags: []*TagValue{},
	}
	p, rslt = newParser(a, "", "")
	api.Sanitize(p)
	rslt.Handler.Stop()
	a.Empty(rslt.Errors)

	api.Tags = append(api.Tags, &TagValue{Content: Content{Value: "s1"}})
	p, rslt = newParser(a, "", "")
	api.Sanitize(p)
	rslt.Handler.Stop()
	a.Empty(rslt.Errors)

	api.Tags = append(api.Tags, &TagValue{Content: Content{Value: "s1"}})
	p, rslt = newParser(a, "", "")
	api.Sanitize(p)
	rslt.Handler.Stop()
	a.NotEmpty(rslt.Errors)
}

func TestXMLnamespace_Sanitize(t *testing.T) {
	a := assert.New(t)

	ns := &XMLNamespace{}
	p, rslt := newParser(a, "", "")
	ns.Sanitize(p)
	rslt.Handler.Stop()
	a.NotEmpty(rslt.Errors)

	ns.URN = &Attribute{Value: xmlenc.String{Value: "urn"}}
	p, rslt = newParser(a, "", "")
	ns.Sanitize(p)
	rslt.Handler.Stop()
	a.Empty(rslt.Errors)
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
			t: &TypeAttribute{Value: xmlenc.String{Value: TypeString}},
			enums: []*Enum{
				{Value: &Attribute{Value: xmlenc.String{Value: "123"}}},
				{Value: &Attribute{Value: xmlenc.String{Value: "str"}}},
			},
		},
		{
			t: &TypeAttribute{Value: xmlenc.String{Value: TypeNumber}},
			enums: []*Enum{
				{Value: &Attribute{Value: xmlenc.String{Value: "123"}}},
				{Value: &Attribute{Value: xmlenc.String{Value: "-1.9"}}},
			},
		},
		{
			t: &TypeAttribute{Value: xmlenc.String{Value: TypeBool}},
			enums: []*Enum{
				{Value: &Attribute{Value: xmlenc.String{Value: "true"}}},
				{Value: &Attribute{Value: xmlenc.String{Value: "0"}}},
			},
		},

		{
			t: &TypeAttribute{Value: xmlenc.String{Value: TypeBool}},
			enums: []*Enum{
				{Value: &Attribute{Value: xmlenc.String{Value: "string"}}},
				{Value: &Attribute{Value: xmlenc.String{Value: "0"}}},
			},
			err: true,
		},
		{
			t: &TypeAttribute{Value: xmlenc.String{Value: TypeNumber}},
			enums: []*Enum{
				{Value: &Attribute{Value: xmlenc.String{Value: "string"}}},
				{Value: &Attribute{Value: xmlenc.String{Value: "0"}}},
			},
			err: true,
		},
		{ // object 是不允许的
			t: &TypeAttribute{Value: xmlenc.String{Value: TypeObject}},
			enums: []*Enum{
				{Value: &Attribute{Value: xmlenc.String{Value: "string"}}},
				{Value: &Attribute{Value: xmlenc.String{Value: "0"}}},
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
		p, rslt := newParser(a, "", "")
		defer rslt.Handler.Stop()
		api.sanitizeTags(p)
	})

	doc := &APIDoc{}
	api.doc = doc
	p, rslt := newParser(a, "", "")
	api.sanitizeTags(p)
	rslt.Handler.Stop()
	a.Empty(rslt.Errors).Empty(rslt.Warns)

	api.Servers = []*ServerValue{
		{Content: Content{Value: "s1"}},
	}
	p, rslt = newParser(a, "", "")
	api.sanitizeTags(p)
	rslt.Handler.Stop()
	a.NotEmpty(rslt.Warns)

	doc.Servers = []*Server{
		{Name: &Attribute{Value: xmlenc.String{Value: "s1"}}},
	}
	p, rslt = newParser(a, "", "")
	api.sanitizeTags(p)
	rslt.Handler.Stop()
	a.Empty(rslt.Errors).Empty(rslt.Warns)

	api.Tags = []*TagValue{
		{Content: Content{Value: "t1"}},
	}
	p, rslt = newParser(a, "", "")
	api.sanitizeTags(p)
	rslt.Handler.Stop()
	a.NotEmpty(rslt.Warns)

	doc.Tags = []*Tag{
		{Name: &Attribute{Value: xmlenc.String{Value: "t1"}}},
	}
	p, rslt = newParser(a, "", "")
	api.sanitizeTags(p)
	rslt.Handler.Stop()
	a.Empty(rslt.Warns)
}

func TestAPI_checkDup(t *testing.T) {
	a := assert.New(t)

	doc := &APIDoc{}
	api := &API{doc: doc}
	p, rslt := newParser(a, "", "")
	api.checkDup(p)
	rslt.Handler.Stop()
	a.Empty(rslt.Errors).Empty(rslt.Warns)

	doc = &APIDoc{
		APIs: []*API{
			{
				Method: &MethodAttribute{Value: xmlenc.String{Value: http.MethodDelete}},
			},
		},
	}
	api = &API{
		doc:    doc,
		Method: &MethodAttribute{Value: xmlenc.String{Value: http.MethodDelete}},
	}
	p, rslt = newParser(a, "", "")
	api.checkDup(p)
	rslt.Handler.Stop()
	a.Empty(rslt.Errors).Empty(rslt.Warns)

	doc = &APIDoc{
		APIs: []*API{
			{
				Method: &MethodAttribute{Value: xmlenc.String{Value: http.MethodDelete}},
			},
			{
				Method: &MethodAttribute{Value: xmlenc.String{Value: http.MethodDelete}},
			},
		},
	}
	api = &API{
		doc:    doc,
		Method: &MethodAttribute{Value: xmlenc.String{Value: http.MethodDelete}},
	}
	p, rslt = newParser(a, "", "")
	api.checkDup(p)
	rslt.Handler.Stop()
	a.NotEmpty(rslt.Errors).Empty(rslt.Warns)

	doc.APIs[0].Servers = []*ServerValue{
		{Content: Content{Value: "s1"}},
		{Content: Content{Value: "s2"}},
	}
	doc.APIs[1].Servers = []*ServerValue{
		{Content: Content{Value: "s1"}},
	}
	api.Servers = []*ServerValue{
		{Content: Content{Value: "s2"}},
	}
	p, rslt = newParser(a, "", "")
	api.checkDup(p)
	rslt.Handler.Stop()
	a.Empty(rslt.Errors).Empty(rslt.Warns)

	api.Servers = []*ServerValue{
		{Content: Content{Value: "s1"}},
	}
	p, rslt = newParser(a, "", "")
	api.checkDup(p)
	rslt.Handler.Stop()
	a.NotEmpty(rslt.Errors).Empty(rslt.Warns)

	// 有空的 path
	doc = &APIDoc{
		APIs: []*API{
			{
				Method: &MethodAttribute{Value: xmlenc.String{Value: http.MethodDelete}},
				Path:   &Path{},
			},
			{
				Method: &MethodAttribute{Value: xmlenc.String{Value: http.MethodDelete}},
			},
		},
	}
	api = &API{
		doc:    doc,
		Method: &MethodAttribute{Value: xmlenc.String{Value: http.MethodDelete}},
	}
	p, rslt = newParser(a, "", "")
	api.checkDup(p)
	rslt.Handler.Stop()
	a.NotEmpty(rslt.Errors).Empty(rslt.Warns)

	// 有非空的 path
	doc = &APIDoc{
		APIs: []*API{
			{
				Method: &MethodAttribute{Value: xmlenc.String{Value: http.MethodDelete}},
				Path:   &Path{Path: &Attribute{Value: xmlenc.String{Value: "/path"}}},
			},
			{
				Method: &MethodAttribute{Value: xmlenc.String{Value: http.MethodDelete}},
			},
		},
	}
	api = &API{
		doc:    doc,
		Method: &MethodAttribute{Value: xmlenc.String{Value: http.MethodDelete}},
	}
	p, rslt = newParser(a, "", "")
	api.checkDup(p)
	rslt.Handler.Stop()
	a.Empty(rslt.Errors).Empty(rslt.Warns)

	doc.APIs[0].Servers = []*ServerValue{
		{Content: Content{Value: "s1"}},
		{Content: Content{Value: "s2"}},
	}
	doc.APIs[1].Servers = []*ServerValue{
		{Content: Content{Value: "s1"}},
	}
	api.Servers = []*ServerValue{
		{Content: Content{Value: "s2"}},
	}
	p, rslt = newParser(a, "", "")
	api.checkDup(p)
	rslt.Handler.Stop()
	a.Empty(rslt.Errors).Empty(rslt.Warns)

	api.Servers = []*ServerValue{
		{Content: Content{Value: "s1"}},
	}
	p, rslt = newParser(a, "", "")
	api.checkDup(p)
	rslt.Handler.Stop()
	a.Empty(rslt.Errors).Empty(rslt.Warns)
}
