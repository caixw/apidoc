// SPDX-License-Identifier: MIT

package ast

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/core/messagetest"
	"github.com/caixw/apidoc/v7/internal/xmlenc"
)

func loadAPIDoc(a *assert.Assertion) *APIDoc {
	data, err := ioutil.ReadFile("./testdata/doc.xml")
	a.NotError(err).NotNil(data)

	doc := &APIDoc{}
	a.NotNil(doc)

	rslt := messagetest.NewMessageHandler()
	doc.Parse(rslt.Handler, core.Block{
		Location: core.Location{
			URI:   "doc.xml",
			Range: core.Range{},
		},
		Data: data,
	})
	rslt.Handler.Stop()
	a.Empty(rslt.Errors)

	// 保证 doc.URI 指向文件地址
	a.NotEmpty(doc.URI.String()).
		True(strings.HasSuffix(doc.URI.String(), "doc.xml"))

	return doc
}

func TestAPIDoc(t *testing.T) {
	a := assert.New(t)
	doc := loadAPIDoc(a)

	a.Equal(doc.BaseTag, xmlenc.BaseTag{
		Base: xmlenc.Base{
			UsageKey: "usage-apidoc",
			Range: core.Range{
				Start: core.Position{Character: 0, Line: 2},
				End:   core.Position{Character: 9, Line: 35},
			},
		},
		StartTag: xmlenc.Name{
			Range: core.Range{
				Start: core.Position{Character: 1, Line: 2},
				End:   core.Position{Character: 7, Line: 2},
			},
			Local: xmlenc.String{
				Range: core.Range{
					Start: core.Position{Character: 1, Line: 2},
					End:   core.Position{Character: 7, Line: 2},
				},
				Value: "apidoc",
			},
		},
		EndTag: xmlenc.Name{
			Range: core.Range{
				Start: core.Position{Character: 2, Line: 35},
				End:   core.Position{Character: 8, Line: 35},
			},
			Local: xmlenc.String{
				Range: core.Range{
					Start: core.Position{Character: 2, Line: 35},
					End:   core.Position{Character: 8, Line: 35},
				},
				Value: "apidoc",
			},
		},
	})

	a.Equal(doc.Version, &VersionAttribute{
		BaseAttribute: xmlenc.BaseAttribute{
			Base: xmlenc.Base{
				UsageKey: "usage-apidoc-version",
				Range: core.Range{
					Start: core.Position{Character: 8, Line: 2},
					End:   core.Position{Character: 23, Line: 2},
				},
			},
			AttributeName: xmlenc.Name{
				Range: core.Range{
					Start: core.Position{Character: 8, Line: 2},
					End:   core.Position{Character: 15, Line: 2},
				},
				Local: xmlenc.String{
					Range: core.Range{
						Start: core.Position{Character: 8, Line: 2},
						End:   core.Position{Character: 15, Line: 2},
					},
					Value: "version",
				},
			},
		},

		Value: xmlenc.String{
			Range: core.Range{
				Start: core.Position{Character: 17, Line: 2},
				End:   core.Position{Character: 22, Line: 2},
			},
			Value: "1.1.1",
		},
	})

	a.Equal(len(doc.Tags), 2)
	tag := &Tag{
		BaseTag: xmlenc.BaseTag{
			Base: xmlenc.Base{
				UsageKey: "usage-apidoc-tags",
				Range: core.Range{
					Start: core.Position{Character: 4, Line: 10},
					End:   core.Position{Character: 47, Line: 10},
				},
			},
			StartTag: xmlenc.Name{
				Range: core.Range{
					Start: core.Position{Character: 5, Line: 10},
					End:   core.Position{Character: 8, Line: 10},
				},
				Local: xmlenc.String{
					Range: core.Range{
						Start: core.Position{Character: 5, Line: 10},
						End:   core.Position{Character: 8, Line: 10},
					},
					Value: "tag",
				},
			},
		},
		Name: &Attribute{
			BaseAttribute: xmlenc.BaseAttribute{
				Base: xmlenc.Base{
					UsageKey: "usage-tag-name",
					Range: core.Range{
						Start: core.Position{Character: 9, Line: 10},
						End:   core.Position{Character: 20, Line: 10},
					},
				},
				AttributeName: xmlenc.Name{
					Range: core.Range{
						Start: core.Position{Character: 9, Line: 10},
						End:   core.Position{Character: 13, Line: 10},
					},
					Local: xmlenc.String{
						Range: core.Range{
							Start: core.Position{Character: 9, Line: 10},
							End:   core.Position{Character: 13, Line: 10},
						},
						Value: "name",
					},
				},
			},
			Value: xmlenc.String{
				Range: core.Range{
					Start: core.Position{Character: 15, Line: 10},
					End:   core.Position{Character: 19, Line: 10},
				},
				Value: "tag1",
			},
		},
		Title: &Attribute{
			BaseAttribute: xmlenc.BaseAttribute{
				Base: xmlenc.Base{
					UsageKey: "usage-tag-title",
					Range: core.Range{
						Start: core.Position{Character: 21, Line: 10},
						End:   core.Position{Character: 44, Line: 10},
					},
				},
				AttributeName: xmlenc.Name{
					Range: core.Range{
						Start: core.Position{Character: 21, Line: 10},
						End:   core.Position{Character: 26, Line: 10},
					},
					Local: xmlenc.String{
						Range: core.Range{
							Start: core.Position{Character: 21, Line: 10},
							End:   core.Position{Character: 26, Line: 10},
						},
						Value: "title",
					},
				},
			},
			Value: xmlenc.String{
				Range: core.Range{
					Start: core.Position{Character: 28, Line: 10},
					End:   core.Position{Character: 43, Line: 10},
				},
				Value: "tag description",
			},
		},
	}
	a.Equal(doc.Tags[0], tag)

	tag = doc.Tags[1]
	a.Equal(tag.Deprecated.V(), "1.0.1").
		Empty(tag.EndTag.Local.Value).
		Equal(tag.UsageKey, "usage-apidoc-tags")

	a.Equal(2, len(doc.Servers))
	srv := doc.Servers[0]
	a.Equal(srv.Name.V(), "admin").
		Equal(srv.URL.V(), "https://api.example.com/admin").
		Nil(srv.Description).
		Equal(srv.Summary.V(), "admin api")

	srv = doc.Servers[1]
	a.Equal(srv.Name.V(), "client").
		Equal(srv.URL.V(), "https://api.example.com/client").
		Equal(srv.Deprecated.V(), "1.0.1").
		Equal(srv.Description.V(), "\n        <p>client api</p>\n        ")

	a.NotNil(doc.License).
		Equal(doc.License.Text.V(), "MIT").
		Equal(doc.License.URL.V(), "https://opensource.org/licenses/MIT")

	a.NotNil(doc.Contact).
		Equal(doc.Contact.Name.V(), "test").
		Equal(doc.Contact.URL.V(), "https://example.com").
		Equal(doc.Contact.Email.V(), "test@example.com")

	a.Contains(doc.Description.V(), "<h2>h2</h2>").
		NotContains(doc.Description.V(), "</description>")

	a.True(doc.tagExists("tag1")).
		False(doc.tagExists("not-exists"))

	a.True(doc.serverExists("admin")).
		False(doc.serverExists("not-exists"))

	a.Equal(2, len(doc.Mimetypes)).
		Equal("application/xml", doc.Mimetypes[0].Content.Value)
}

func TestAPIDoc_all(t *testing.T) {
	a := assert.New(t)

	data, err := ioutil.ReadFile("./testdata/all.xml")
	a.NotError(err).NotNil(data)
	rslt := messagetest.NewMessageHandler()
	doc := &APIDoc{}
	doc.Parse(rslt.Handler, core.Block{Data: data, Location: core.Location{URI: "all.xml"}})
	rslt.Handler.Stop()
	a.Empty(rslt.Errors)

	a.Equal(doc.Version.V(), "1.1.1").False(doc.Version.AttributeName.IsEmpty())

	a.Equal(len(doc.Tags), 2)
	tag := doc.Tags[0]
	a.Equal(tag.Name.V(), "tag1").
		NotEmpty(tag.Title.V())
	tag = doc.Tags[1]
	a.Equal(tag.Deprecated.V(), "1.0.1").
		Equal(tag.Name.V(), "tag2")

	a.Equal(2, len(doc.Servers))
	srv := doc.Servers[0]
	a.Equal(srv.Name.V(), "admin").
		Equal(srv.URL.V(), "https://api.example.com/admin").
		Nil(srv.Description)

	a.True(doc.tagExists("tag1")).
		False(doc.tagExists("not-exists"))

	a.True(doc.serverExists("admin")).
		False(doc.serverExists("not-exists"))

	a.Equal(2, len(doc.Mimetypes)).
		Equal("application/xml", doc.Mimetypes[0].Content.Value)

	// api
	a.Equal(1, len(doc.APIs))

	a.NotEmpty(doc.URI).
		Empty(doc.APIs[0].URI.String())
}

func loadAPI(a *assert.Assertion) *API {
	doc := loadAPIDoc(a)

	data, err := ioutil.ReadFile("./testdata/api.xml")
	a.NotError(err).NotNil(data)

	rslt := messagetest.NewMessageHandler()
	doc.Parse(rslt.Handler, core.Block{
		Data:     data,
		Location: core.Location{URI: core.FileURI("./testdata/api.xml")},
	})
	rslt.Handler.Stop()
	a.Empty(rslt.Errors)

	api := doc.APIs[0]

	// 保证 doc.APIs[0].URI 指向文件地址
	a.NotEmpty(api.URI.String()).
		True(strings.HasSuffix(api.URI.String(), "./testdata/api.xml"))

	return api
}

func TestAPI(t *testing.T) {
	a := assert.New(t)
	api := loadAPI(a)

	a.Equal(api.Version.V(), "1.1.0").
		Equal(2, len(api.Tags))

	a.Equal(len(api.Responses), 2)
	resp := api.Responses[0]
	a.Equal(resp.Mimetype.V(), "json").
		Equal(resp.Status.V(), 200).
		Equal(resp.Type.V(), TypeObject).
		Equal(len(resp.Items), 3)
	sex := resp.Items[1]
	a.Equal(sex.Type.V(), TypeString).
		Equal(sex.Default.V(), "male").
		Equal(len(sex.Enums), 2)
	example := resp.Examples[0]
	a.Equal(example.Mimetype.V(), "json").
		NotEmpty(example.Content.Value.Value)

	a.Equal(1, len(api.Requests))
	req := api.Requests[0]
	a.Equal(req.Mimetype.V(), "json").
		Equal(req.Headers[0].Name.V(), "authorization")

	// callback
	cb := api.Callback
	a.Equal(cb.Method.V(), "POST").
		Equal(cb.Requests[0].Type.V(), TypeObject).
		Equal(cb.Requests[0].Mimetype.V(), "json").
		Equal(cb.Responses[0].Status.V(), 200)
}

func TestRequest_Param(t *testing.T) {
	a := assert.New(t)

	var req *Request
	a.Nil(req.Param())

	req = &Request{Type: &TypeAttribute{Value: xmlenc.String{Value: TypeObject}}}
	param := req.Param()
	a.Equal(req.Type, param.Type)
}

func TestAPIDoc_XMLNamespaces(t *testing.T) {
	a := assert.New(t)

	d := &APIDoc{
		XMLNamespaces: []*XMLNamespace{
			{
				URN: &Attribute{Value: xmlenc.String{Value: core.XMLNamespace}},
			},
			{
				URN:    &Attribute{Value: xmlenc.String{Value: "urn1"}},
				Prefix: &Attribute{Value: xmlenc.String{Value: "ns1"}},
			},
			{
				URN:    &Attribute{Value: xmlenc.String{Value: "urn2"}},
				Prefix: &Attribute{Value: xmlenc.String{Value: "ns2"}},
			},
		},
	}

	ns := d.XMLNamespace("")
	a.Empty(ns.Prefix.V()).Equal(ns.URN.V(), core.XMLNamespace)

	ns = d.XMLNamespace("ns1")
	a.Equal(ns.URN.V(), "urn1")

	ns = d.XMLNamespace("ns2")
	a.Equal(ns.URN.V(), "urn2")

	// not exists
	ns = d.XMLNamespace("ns")
	a.Nil(ns)
}
