// Copyright 2019 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

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
	a.Equal(api.Version, "1.1").
		Equal(api.Tags, []string{"g1", "g2"}).
		Equal(len(api.Responses), 2)

	resp := api.Responses[0]
	a.Equal(resp.Mimetype, "json,xml").
		Equal(resp.Status, 200).
		Equal(resp.Type, Object).
		Equal(len(resp.Items), 3)
	sex := resp.Items[1]
	a.Equal(sex.Type, String).
		Equal(sex.Default, "male").
		Equal(len(sex.Enums), 2)
}
