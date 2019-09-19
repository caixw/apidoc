// SPDX-License-Identifier: MIT

package input

import (
	"bytes"
	"testing"

	"github.com/issue9/assert"
)

var (
	api1 = []byte(`@api POST /users/login 登录
group users
tags: [t1,t2]

request:
  description: request body
  content:
    application/json:
      schema:
        type: object
        properties:
          username:
            type: string
            description: 登录账号
          password:
            type: string
            description: 密码
`)

	api2 = []byte(`@api DELETE /users/login 注销登录
group users
tags: [t1,t2]

request:
  description: request body
  content:
    application/json:
      schema:
        type: object
        properties:
          username:
            type: string
            description: 登录账号
          password:
            type: string
            description: 密码
`)

	doc1 = []byte(`@apidoc title of api
version: 2.9
license:
  name: MIT
  url: https://opensources.org/licenses/MIT
description:>
  line1
  line2
`)
)

func TestParse(t *testing.T) {
	a := assert.New(t)

	testBuildBlock(a, "go")
	testBuildBlock(a, "groovy")
	testBuildBlock(a, "java")
	testBuildBlock(a, "javascript")
	testBuildBlock(a, "pascal")
	testBuildBlock(a, "perl")
	testBuildBlock(a, "php")
	testBuildBlock(a, "python")
	testBuildBlock(a, "ruby")
	testBuildBlock(a, "rust")
	testBuildBlock(a, "swift")
}

func testBuildBlock(a *assert.Assertion, lang string) {
	o := &Options{
		Lang:      lang,
		Dir:       "./testdata/" + lang,
		Recursive: true,
	}
	a.NotError(o.Sanitize())

	channel := buildBlock(nil, o)
	a.NotNil(channel)

	for b := range channel {
		eq := bytes.Equal(b.Data, api1) ||
			bytes.Equal(b.Data, api2) ||
			bytes.Equal(b.Data, doc1) ||
			(!bytes.HasPrefix(b.Data, []byte("@api ")) && !bytes.HasPrefix(b.Data, []byte("@apidoc ")))
		a.True(eq, "lang(%s)：%s,%s,%d,%d", lang, string(b.Data), string(api1), len(b.Data), len(api1))
	}
}
