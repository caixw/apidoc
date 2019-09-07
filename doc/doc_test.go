// SPDX-License-Identifier: MIT

package doc

import (
	"encoding/xml"
	"io/ioutil"
	"testing"

	"github.com/issue9/assert"
)

func TestDoc(t *testing.T) {
	a := assert.New(t)

	data, err := ioutil.ReadFile("./doc.xml")
	a.NotError(err).NotNil(data)

	doc := &Doc{}
	a.NotError(xml.Unmarshal(data, doc))
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
}
