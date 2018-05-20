// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package output

import (
	"encoding/json"
	"testing"

	"github.com/issue9/assert"
	yaml "gopkg.in/yaml.v2"
)

func TestOptions_Sanitize(t *testing.T) {
	a := assert.New(t)
	o := &Options{}
	a.Error(o.Sanitize())

	o.Dir = "./testdir"
	a.NotError(o.Sanitize())
	a.Equal(o.marshal, yaml.Marshal)

	o.Type = typeJSON
	a.NotError(o.Sanitize())
	a.Equal(o.marshal, json.Marshal)

	o.Type = "unknown"
	a.Error(o.Sanitize())

	// 会执行删除 testdir 操作
	o.Type = typeJSON
	o.Clean = true
	a.NotError(o.Sanitize())
}

func TestOptions_contains(t *testing.T) {
	a := assert.New(t)

	o := &Options{}
	a.True(o.contains("not exists"))

	o.Groups = []string{"g1", "g2"}
	a.True(o.contains("g1")).
		True(o.contains("g2")).
		False(o.contains("not exists"))
}
