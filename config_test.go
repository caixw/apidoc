// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/issue9/assert"
)

func TestGenConfigFile(t *testing.T) {
	a := assert.New(t)
	a.NotError(genConfigFile())

	data, err := ioutil.ReadFile("./" + configFilename)
	a.NotError(err).NotNil(data)
	cfg := &config{}
	a.NotError(json.Unmarshal(data, cfg))

	a.Equal(cfg.Input.Dir, "./").Equal(cfg.Input.Recursive, true)
}
