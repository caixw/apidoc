// SPDX-License-Identifier: MIT

package doc

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v5/message"
)

func loadDoc(a *assert.Assertion) *Doc {
	data, err := ioutil.ReadFile("./testdata/doc.xml")
	a.NotError(err).NotNil(data)

	doc := New()
	a.NotNil(doc).NotEmpty(doc.APIDoc)

	a.NotError(doc.FromXML("doc.xml", 0, data))

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
}

func TestDoc_all(t *testing.T) {
	a := assert.New(t)

	data, err := ioutil.ReadFile("./testdata/all.xml")
	a.NotError(err).NotNil(data)
	doc := New()
	a.NotNil(doc)
	a.NotError(doc.FromXML("all.xml", 0, data))

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

	// api
	a.Equal(1, len(doc.Apis))
}

func TestDoc_UnmarshalXML(t *testing.T) {
	a := assert.New(t)

	// 重得的标签名
	data := `<apidoc version="1.1.1">
		<tag name="t1" title="tet" />
		<tag name="t1" title="tet"></tag>
	</apidoc>`
	doc := New()
	a.NotNil(doc)
	err := doc.FromXML("file", 11, []byte(data))
	serr, ok := err.(*message.SyntaxError)
	a.True(ok).NotNil(serr)
	a.Equal(serr.Line, 11).
		Equal(serr.File, "file")

		// 重复得的 server
	data = `<apidoc version="1.1.1">
		<server name="s1" url="https://example.com/s1" summary="tet" />
		<server name="s1" url="https://example.com/s2" summary="tet" />
	</apidoc>`
	doc = New()
	a.NotNil(doc)
	err = doc.FromXML("file", 12, []byte(data))
	serr, ok = err.(*message.SyntaxError)
	a.True(ok).NotNil(serr)
	a.Equal(serr.Line, 12).
		Equal(serr.File, "file")

	data = `<apidoc version="1.1.1">
			<tag name="t1" deprecated="x.0.1" />
		</apidoc>`
	err = doc.FromXML("file", 11, []byte(data))
	serr, ok = err.(*message.SyntaxError)
	a.True(ok).NotNil(serr)
	a.Equal(serr.Line, 12)
}

func TestDoc_Sanitize(t *testing.T) {
	a := assert.New(t)
	doc := New()
	a.NotNil(doc)

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
		},
		{
			Tags:   []string{"not-exists", "tag1"},
			doc:    doc,
			Path:   &Path{},
			Method: http.MethodGet,
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
		},
		{
			Servers: []string{"not-exists", "tag1"},
			doc:     doc,
			Path:    &Path{},
			Method:  http.MethodGet,
		},
	}
	a.Error(doc.Sanitize())
}

// 测试错误提示的行号是否正确
func TestDoc_lineNumber(t *testing.T) {
	a := assert.New(t)
	doc := New()
	a.NotNil(doc)

	data := []byte(`<apidoc version="x.0.1"></apidoc>`)
	err := doc.FromXML("file", 11, data)
	a.Equal(err.(*message.SyntaxError).Line, 11)

	data = []byte(`<apidoc
	
	version="x.1.1">
	</apidoc>`)
	err = doc.FromXML("file", 12, data)
	a.Equal(err.(*message.SyntaxError).Line, 14)
}
