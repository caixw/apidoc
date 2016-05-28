// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package input

import (
	"testing"

	"github.com/caixw/apidoc/doc"
	"github.com/issue9/assert"
)

func TestParseFile(t *testing.T) {
	a := assert.New(t)

	testParseFile(a, "go", "./testdata/go/test1.go")
	testParseFile(a, "php", "./testdata/php/test1.php")
	testParseFile(a, "c", "./testdata/c/test1.c")
	testParseFile(a, "ruby", "./testdata/ruby/test1.rb")
}

func testParseFile(a *assert.Assertion, lang string, path string) {
	docs := doc.New()
	a.NotNil(docs)

	b, found := langs[lang]
	if !found {
		a.T().Error("不支持该语言")
	}

	parseFile(docs, path, b)
	a.Equal(2, len(docs.Apis))

	api0 := docs.Apis[0]
	api1 := docs.Apis[1]

	a.Equal(api0.URL, "/users/login").
		Equal(api1.URL, "/users/login").
		Equal(api0.Group, "users").
		Equal(api1.Group, "users")

	if api0.Method == "POST" {
		a.Equal(api1.Method, "DELETE").
			Equal(1, len(api1.Request.Headers))
	} else {
		a.Equal(api0.Method, "DELETE").
			Equal(api1.Method, "POST").
			Equal(1, len(api0.Request.Headers))
	}
}
