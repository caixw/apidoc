// SPDX-License-Identifier: MIT

package spec

import (
	"io/ioutil"
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v6/core"
)

func loadAPI(a *assert.Assertion) *API {
	doc := loadDoc(a)

	data, err := ioutil.ReadFile("./testdata/api.xml")
	a.NotError(err).NotNil(data)

	a.NotError(doc.appendAPI(&core.Block{Data: data}))
	return doc.Apis[0]
}

func TestAPI(t *testing.T) {
	a := assert.New(t)
	api := loadAPI(a)

	a.Equal(api.Version, "1.1.0").
		Equal(api.Tags, []string{"g1", "g2"})

	a.Equal(len(api.Responses), 2)
	resp := api.Responses[0]
	a.Equal(resp.Mimetype, "json").
		Equal(resp.Status, 200).
		Equal(resp.Type, Object).
		Equal(len(resp.Items), 3)
	sex := resp.Items[1]
	a.Equal(sex.Type, String).
		Equal(sex.Default, "male").
		Equal(len(sex.Enums), 2)
	example := resp.Examples[0]
	a.Equal(example.Mimetype, "json").
		NotEmpty(example.Content)

	a.Equal(1, len(api.Requests))
	req := api.Requests[0]
	a.Equal(req.Mimetype, "json").
		Equal(req.Headers[0].Name, "authorization")

	// callback
	cb := api.Callback
	a.Equal(cb.Method, "POST").
		Equal(cb.Requests[0].Type, Object).
		Equal(cb.Requests[0].Mimetype, "json").
		Equal(cb.Responses[0].Status, 200)
}

func TestAPI_UnmarshalXML(t *testing.T) {
	a := assert.New(t)

	doc := NewAPIDoc()
	data := `<api version="1.1.1">
		<header type="object" name="key1" summary="summary">
			<param name="id" type="number" summary="summary" />
		</header>
		<response type="number" summary="summary" />
	</api>`
	a.Error(doc.appendAPI(&core.Block{Data: []byte(data)}))
}
