// SPDX-License-Identifier: MIT

package ast

import (
	"io/ioutil"
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v6/core"
	"github.com/caixw/apidoc/v6/internal/token"
)

func loadAPIDoc(a *assert.Assertion) *APIDoc {
	data, err := ioutil.ReadFile("./testdata/doc.xml")
	a.NotError(err).NotNil(data)

	doc := &APIDoc{}
	a.NotNil(doc)

	a.NotError(doc.Parse(core.Block{
		Location: core.Location{
			URI:   "doc.xml",
			Range: core.Range{},
		},
		Data: data,
	}))

	return doc
}

func TestAPIDoc(t *testing.T) {
	a := assert.New(t)
	doc := loadAPIDoc(a)

	a.Equal(doc.Version, &VersionAttribute{
		Base: token.Base{
			UsageKey: "usage-apidoc-version",
			Range: core.Range{
				Start: core.Position{Character: 8, Line: 2},
				End:   core.Position{Character: 23, Line: 2},
			},
			XMLName: String{
				Range: core.Range{
					Start: core.Position{Character: 8, Line: 2},
					End:   core.Position{Character: 15, Line: 2},
				},
				Value: "version",
			},
		},

		Value: String{
			Range: core.Range{
				Start: core.Position{Character: 17, Line: 2},
				End:   core.Position{Character: 22, Line: 2},
			},
			Value: "1.1.1",
		},
	})

	a.Equal(len(doc.Tags), 2)
	tag := &Tag{
		Base: token.Base{
			UsageKey: "usage-apidoc-tags",
			Range: core.Range{
				Start: core.Position{Character: 4, Line: 10},
				End:   core.Position{Character: 47, Line: 10},
			},
			XMLName: String{
				Range: core.Range{
					Start: core.Position{Character: 5, Line: 10},
					End:   core.Position{Character: 8, Line: 10},
				},
				Value: "tag",
			},
		},
		Name: &Attribute{
			Base: token.Base{
				UsageKey: "usage-tag-name",
				Range: core.Range{
					Start: core.Position{Character: 9, Line: 10},
					End:   core.Position{Character: 20, Line: 10},
				},
				XMLName: String{
					Range: core.Range{
						Start: core.Position{Character: 9, Line: 10},
						End:   core.Position{Character: 13, Line: 10},
					},
					Value: "name",
				},
			},
			Value: String{
				Range: core.Range{
					Start: core.Position{Character: 15, Line: 10},
					End:   core.Position{Character: 19, Line: 10},
				},
				Value: "tag1",
			},
		},
		Title: &Attribute{
			Base: token.Base{
				UsageKey: "usage-tag-title",
				Range: core.Range{
					Start: core.Position{Character: 21, Line: 10},
					End:   core.Position{Character: 44, Line: 10},
				},
				XMLName: String{
					Range: core.Range{
						Start: core.Position{Character: 21, Line: 10},
						End:   core.Position{Character: 26, Line: 10},
					},
					Value: "title",
				},
			},
			Value: String{
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
	a.Equal(tag.Deprecated.Value.Value, "1.0.1").
		Empty(tag.XMLNameEnd.Value).
		Equal(tag.UsageKey, "usage-apidoc-tags")

	a.Equal(2, len(doc.Servers))
	srv := doc.Servers[0]
	a.Equal(srv.Name.Value.Value, "admin").
		Equal(srv.URL.Value.Value, "https://api.example.com/admin").
		Nil(srv.Description).
		Equal(srv.Summary.Value.Value, "admin api")

	srv = doc.Servers[1]
	a.Equal(srv.Name.Value.Value, "client").
		Equal(srv.URL.Value.Value, "https://api.example.com/client").
		Equal(srv.Deprecated.Value.Value, "1.0.1").
		Equal(srv.Description.Text.Value.Value, "\n        <p>client api</p>\n        ")

	a.NotNil(doc.License).
		Equal(doc.License.Text.Value.Value, "MIT").
		Equal(doc.License.URL.Value.Value, "https://opensource.org/licenses/MIT")

	a.NotNil(doc.Contact).
		Equal(doc.Contact.Name.Value.Value, "test").
		Equal(doc.Contact.URL.Content.Value, "https://example.com").
		Equal(doc.Contact.Email.Content.Value, "test@example.com")

	a.NotEmpty(doc.Description.Text).
		Contains(doc.Description.Text.Value.Value, "<h2>h2</h2>").
		NotContains(doc.Description.Text.Value.Value, "</description>")

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
	doc := &APIDoc{}
	a.NotError(doc.Parse(core.Block{Data: data}))

	a.Equal(doc.Version.Value.Value, "1.1.1").False(doc.Version.XMLName.IsEmpty())

	a.Equal(len(doc.Tags), 2)
	tag := doc.Tags[0]
	a.Equal(tag.Name.Value.Value, "tag1").
		NotEmpty(tag.Title.Value.Value)
	tag = doc.Tags[1]
	a.Equal(tag.Deprecated.Value.Value, "1.0.1").
		Equal(tag.Name.Value.Value, "tag2")

	a.Equal(2, len(doc.Servers))
	srv := doc.Servers[0]
	a.Equal(srv.Name.Value.Value, "admin").
		Equal(srv.URL.Value.Value, "https://api.example.com/admin").
		Nil(srv.Description)

	a.True(doc.tagExists("tag1")).
		False(doc.tagExists("not-exists"))

	a.True(doc.serverExists("admin")).
		False(doc.serverExists("not-exists"))

	a.Equal(2, len(doc.Mimetypes)).
		Equal("application/xml", doc.Mimetypes[0].Content.Value)

	// api
	a.Equal(1, len(doc.Apis))
}

func loadAPI(a *assert.Assertion) *API {
	doc := loadAPIDoc(a)

	data, err := ioutil.ReadFile("./testdata/api.xml")
	a.NotError(err).NotNil(data)

	a.NotError(doc.Parse(core.Block{Data: data}))
	return doc.Apis[0]
}

func TestAPI(t *testing.T) {
	a := assert.New(t)
	api := loadAPI(a)

	a.Equal(api.Version.Value.Value, "1.1.0").
		Equal(2, len(api.Tags))

	a.Equal(len(api.Responses), 2)
	resp := api.Responses[0]
	a.Equal(resp.Mimetype.Value.Value, "json").
		Equal(resp.Status.Value.Value, 200).
		Equal(resp.Type.Value.Value, TypeObject).
		Equal(len(resp.Items), 3)
	sex := resp.Items[1]
	a.Equal(sex.Type.Value.Value, TypeString).
		Equal(sex.Default.Value.Value, "male").
		Equal(len(sex.Enums), 2)
	example := resp.Examples[0]
	a.Equal(example.Mimetype.Value.Value, "json").
		NotEmpty(example.Content.Value.Value)

	a.Equal(1, len(api.Requests))
	req := api.Requests[0]
	a.Equal(req.Mimetype.Value.Value, "json").
		Equal(req.Headers[0].Name.Value.Value, "authorization")

	// callback
	cb := api.Callback
	a.Equal(cb.Method.Value.Value, "post").
		Equal(cb.Requests[0].Type.Value.Value, TypeObject).
		Equal(cb.Requests[0].Mimetype.Value.Value, "json").
		Equal(cb.Responses[0].Status.Value.Value, 200)
}
