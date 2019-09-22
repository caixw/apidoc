// SPDX-License-Identifier: MIT

package doc

import (
	"io/ioutil"
	"testing"

	"github.com/issue9/assert"
)

func newAPI(a *assert.Assertion) *API {
	doc := newDoc(a)

	data, err := ioutil.ReadFile("./testdata/api.xml")
	a.NotError(err).NotNil(data)

	api := doc.NewAPI("", 1)
	a.NotNil(api)

	a.NotError(api.FromXML(data))

	return api
}

func TestAPI(t *testing.T) {
	a := assert.New(t)
	api := newAPI(a)

	a.Equal(api.Version, "1.1.0").
		Equal(api.Tags, []string{"g1", "g2"})

	a.Equal(len(api.Responses), 2)
	resp := api.Responses[0]
	a.Equal(resp.Mimetype, "json,xml").
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
	a.Equal(req.Mimetype, "json,xml").
		Equal(req.Headers[0].Name, "authorization")

	// callback
	cb := api.Callback
	a.Equal(cb.Method, "POST").
		Equal(cb.Schema, "https").
		Equal(cb.Requests[0].Type, Object).
		Equal(cb.Requests[0].Mimetype, "json").
		Equal(cb.Responses[0].Status, 200)
}
