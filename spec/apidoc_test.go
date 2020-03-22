// SPDX-License-Identifier: MIT

package spec

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v6/core"
)

func loadDoc(a *assert.Assertion) *APIDoc {
	data, err := ioutil.ReadFile("./testdata/doc.xml")
	a.NotError(err).NotNil(data)

	doc := NewAPIDoc()
	a.NotNil(doc).NotEmpty(doc.APIDoc)

	a.NotError(doc.fromXML(&Block{
		Location: core.Location{
			URI:   "doc.xml",
			Range: core.Range{},
		},
		Data: data,
	}))

	return doc
}

func TestDoc(t *testing.T) {
	a := assert.New(t)
	doc := loadDoc(a)

	a.NotEmpty(doc.APIDoc)
	a.Equal(doc.Version, "1.1.1")

	a.Equal(len(doc.Tags), 2)
	tag := doc.Tags[0]
	a.Equal(tag.Name, "tag1").NotEmpty(tag.Title)
	tag = doc.Tags[1]
	a.Equal(tag.Deprecated, "1.0.1").Equal(tag.Name, "tag2")

	a.Equal(2, len(doc.Servers))
	srv := doc.Servers[0]
	a.Equal(srv.Name, "admin").
		Equal(srv.URL, "https://api.example.com/admin").
		NotEmpty(srv.Description)
	srv = doc.Servers[1]
	a.Equal(srv.Name, "client").
		Equal(srv.URL, "https://api.example.com/client").
		Equal(srv.Deprecated, "1.0.1")

	a.Equal(doc.License, &Link{Text: "MIT", URL: "https://opensource.org/licenses/MIT"}).
		Equal(doc.Contact, &Contact{
			Name:  "test",
			URL:   "https://example.com",
			Email: "test@example.com",
		})

	a.NotEmpty(doc.Description.Text).
		Contains(doc.Description.Text, "<h2>h2</h2>").
		NotContains(doc.Description.Text, "</description>")

	a.True(doc.tagExists("tag1")).
		False(doc.tagExists("not-exists"))

	a.True(doc.serverExists("admin")).
		False(doc.serverExists("not-exists"))

	a.Equal(2, len(doc.Mimetypes)).
		Equal("application/xml", doc.Mimetypes[0])
}

func TestDoc_all(t *testing.T) {
	a := assert.New(t)

	data, err := ioutil.ReadFile("./testdata/all.xml")
	a.NotError(err).NotNil(data)
	doc := NewAPIDoc()
	a.NotNil(doc)
	a.NotError(doc.fromXML(&Block{Data: data}))

	a.Equal(doc.Version, "1.1.1")

	a.Equal(len(doc.Tags), 2)
	tag := doc.Tags[0]
	a.Equal(tag.Name, "tag1").NotEmpty(tag.Title)
	tag = doc.Tags[1]
	a.Equal(tag.Deprecated, "1.0.1").Equal(tag.Name, "tag2")

	a.Equal(2, len(doc.Servers))
	srv := doc.Servers[0]
	a.Equal(srv.Name, "admin").
		Equal(srv.URL, "https://api.example.com/admin").
		NotEmpty(srv.Description)

	a.True(doc.tagExists("tag1")).
		False(doc.tagExists("not-exists"))

	a.True(doc.serverExists("admin")).
		False(doc.serverExists("not-exists"))

	a.Equal(2, len(doc.Mimetypes)).
		Equal("application/xml", doc.Mimetypes[0])

	// api
	a.Equal(1, len(doc.Apis))
}

func TestDoc_UnmarshalXML(t *testing.T) {
	a := assert.New(t)
	loc := core.Location{
		URI: "file:///file.php",
		Range: core.Range{
			Start: core.Position{
				Line:      11,
				Character: 22,
			},
			End: core.Position{},
		},
	}

	// 重得的标签名
	data := `<apidoc version="1.1.1">
		<tag name="t1" title="tet" />
		<tag name="t1" title="tet"></tag>
		<mimetype>application/json</mimetype>
		<title>title</title>
	</apidoc>`
	doc := NewAPIDoc()
	a.NotNil(doc)
	err := doc.fromXML(&Block{Location: loc, Data: []byte(data)})
	serr, ok := err.(*core.SyntaxError)
	a.True(ok).NotNil(serr)
	a.Equal(serr.Location.Range.Start.Line, 11).
		Equal(serr.Location.URI, "file:///file.php")

	// 缺少 title
	doc = NewAPIDoc()
	data = `<apidoc version="1.1.1">
			<mimetype>application/json</mimetype>
	</apidoc>`
	err = doc.fromXML(&Block{Location: loc, Data: []byte(data)})
	serr, ok = err.(*core.SyntaxError)
	a.True(ok).NotNil(serr)
	a.Equal(serr.Location.Range.Start.Line, 11)

	// 重复得的 server
	doc = NewAPIDoc()
	data = `<apidoc version="1.1.1">
		<server name="s1" url="https://example.com/s1" summary="tet" />
		<server name="s1" url="https://example.com/s2" summary="tet" />
		<mimetype>application/json</mimetype>
		<title>title</title>
	</apidoc>`
	a.NotNil(doc)
	err = doc.fromXML(&Block{Location: loc, Data: []byte(data)})
	serr, ok = err.(*core.SyntaxError)
	a.True(ok).NotNil(serr)
	a.Equal(serr.Location.Range.Start.Line, 11).
		Equal(serr.Location.URI, "file:///file.php")

	// 无效的 deprecated 值
	doc = NewAPIDoc()
	data = `<apidoc version="1.1.1">
			<tag name="t1" deprecated="x.0.1" />
			<mimetype>application/json</mimetype>
			<title>title</title>
		</apidoc>`
	err = doc.fromXML(&Block{Location: loc, Data: []byte(data)})
	serr, ok = err.(*core.SyntaxError)
	a.True(ok).NotNil(serr)
	a.Equal(serr.Location.Range.Start.Line, 12)

	// 缺少 mimetype
	doc = NewAPIDoc()
	data = `<apidoc version="1.1.1">
			<title>title</title>
	</apidoc>`
	err = doc.fromXML(&Block{Location: loc, Data: []byte(data)})
	serr, ok = err.(*core.SyntaxError)
	a.True(ok).NotNil(serr)
	a.Equal(serr.Location.Range.Start.Line, 11)

	// response.header.type 错误
	doc = NewAPIDoc()
	data = `<apidoc version="1.1.1">
		<mimetype>application/json</mimetype>
		<title>title</title>
		<response type="number" summary="summary">
			<header type="object" name="key1" summary="summary">
				<param name="id" type="number" summary="summary" />
			</header>
		</response>
	</apidoc>`
	a.Error(doc.fromXML(&Block{Location: loc, Data: []byte(data)}))
}

func TestDoc_Sanitize(t *testing.T) {
	a := assert.New(t)
	doc := NewAPIDoc()
	a.NotNil(doc)
	doc.Block = &Block{}

	// api.tags 不存在于 doc
	doc.Tags = []*Tag{
		{Name: "tag1"},
		{Name: "tag2"},
	}
	doc.Apis = []*API{
		{
			Tags:   []string{"tag1", "tag2"},
			doc:    doc,
			Path:   &Path{},
			Method: http.MethodGet,
			Block:  &Block{},
		},
		{
			Tags:   []string{"not-exists", "tag1"},
			doc:    doc,
			Path:   &Path{},
			Method: http.MethodGet,
			Block:  &Block{},
		},
	}
	a.Error(doc.Sanitize())

	// api.servers 不存在于 doc
	doc.Servers = []*Server{
		{Name: "tag1"},
		{Name: "tag2"},
	}
	doc.Apis = []*API{
		{
			Servers: []string{"tag1", "tag2"},
			doc:     doc,
			Path:    &Path{},
			Method:  http.MethodGet,
			Block:   &Block{},
		},
		{
			Servers: []string{"not-exists", "tag1"},
			doc:     doc,
			Path:    &Path{},
			Method:  http.MethodGet,
			Block:   &Block{},
		},
	}
	a.Error(doc.Sanitize())
}
