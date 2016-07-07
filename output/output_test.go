// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package output

import (
	"testing"

	"github.com/issue9/assert"
)

func TestOptions_Init(t *testing.T) {
	a := assert.New(t)

	o := &Options{Type: "html"}
	a.Error(o.Init()) // 未指定 dir

	o.Dir = "./"
	a.NotError(o.Init())

	// 模板不存在
	o.Template = "./not_exists"
	a.Error(o.Init())

	// 未指定 template
	o.Type = "html+"
	o.Template = ""
	a.Error(o.Init())

	// port 未指定
	o.Type = "html+"
	o.Template = "./static"
	a.Error(o.Init())

	// 修正 port
	o.Port = "1234"
	a.NotError(o.Init())
	a.Equal(":1234", o.Port)
}

func TestIsSupportedType(t *testing.T) {
	a := assert.New(t)

	a.True(isSuppertedType("html"))
	a.True(isSuppertedType("html+"))
	a.True(isSuppertedType("json"))
	a.False(isSuppertedType("not-exists"))
}
