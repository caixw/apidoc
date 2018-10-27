// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package doc

import (
	"bytes"
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/doc/lexer"
)

func TestBody_parseExample(t *testing.T) {
	a := assert.New(t)
	body := &Body{}

	a.NotError(body.parseExample(&lexer.Tag{Data: []byte(`application/json summary text
{
	"id": 1,
	"name": "name"
}`)}))
	e := body.Examples[0]
	a.Equal(e.Mimetype, "application/json").
		Equal(e.Summary, "summary text").
		Equal(e.Value, `{
	"id": 1,
	"name": "name"
}`)

	// 长度不够
	a.Error(body.parseExample(&lexer.Tag{Data: []byte("application/json")}))
}

func TestBody_parseHeader(t *testing.T) {
	a := assert.New(t)
	body := &Body{}

	a.NotError(body.parseHeader(&lexer.Tag{Data: []byte(`content-type required json 或是 xml`)}))
	h := body.Headers[0]
	a.Equal(h.Summary, "json 或是 xml").
		Equal(h.Name, "content-type").
		False(h.Optional)

	a.NotError(body.parseHeader(&lexer.Tag{Data: []byte(`ETag optional etag`)}))
	h = body.Headers[1]
	a.Equal(h.Summary, "etag").
		Equal(h.Name, "ETag").
		True(h.Optional)

	// 长度不够
	a.Error(body.parseHeader(&lexer.Tag{Data: []byte("ETag")}))
}

func TestIsOptional(t *testing.T) {
	a := assert.New(t)

	a.False(isOptional(requiredBytes))
	a.False(isOptional(bytes.ToUpper(requiredBytes)))
	a.True(isOptional([]byte("optional")))
	a.True(isOptional([]byte("Optional")))
}
