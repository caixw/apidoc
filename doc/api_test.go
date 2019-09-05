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
	// TODO
}
