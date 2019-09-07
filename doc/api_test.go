// SPDX-License-Identifier: MIT

package doc

import (
	"encoding/xml"
	"io/ioutil"
	"testing"

	"github.com/issue9/assert"
)

func TestAPI(t *testing.T) {
	a := assert.New(t)

	data, err := ioutil.ReadFile("./api.xml")
	a.NotError(err).NotNil(data)

	api := &API{}
	a.NotError(xml.Unmarshal(data, api))
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
	a.Equal(cb.Method, "post").
		Equal(cb.Mimetype, "json").
		Equal(cb.Schema, "https").
		Equal(cb.Type, Object).
		Equal(cb.Responses[0].Status, 200)
}
