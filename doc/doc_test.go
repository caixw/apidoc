// SPDX-License-Identifier: MIT

package doc

import (
	"io/ioutil"
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v5/message"
)

func newDoc(a *assert.Assertion) *Doc {
	data, err := ioutil.ReadFile("./testdata/doc.xml")
	a.NotError(err).NotNil(data)

	doc := New()
	a.NotNil(doc)

	a.NotError(doc.FromXML("doc.xml", 0, data))

	return doc
}

func TestDoc(t *testing.T) {
	a := assert.New(t)
	doc := newDoc(a)

	a.Equal(doc.Version, "1.1.1")

	a.Equal(len(doc.Tags), 2)
	tag := doc.Tags[0]
	a.Equal(tag.Name, "tag1").NotEmpty(tag.Description)
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

	a.True(doc.tagExists("tag1")).
		False(doc.tagExists("not-exists"))

	a.True(doc.serverExists("admin")).
		False(doc.serverExists("not-exists"))
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
