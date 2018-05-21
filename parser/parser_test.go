// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package parser

import (
	"testing"

	"github.com/issue9/assert"
)

var (
	docContent = []byte(`@apidoc title
group: g1
servers:
  - url: https://example.com/{version}
    description: 版本号
    variables:
      - version:
          enum: [v1,v2]
          default: v1
tags:
  - name: tag1
    description: tag1 desc
  - name: tag2
    description: tag2 desc
version: 2.0
description: | ## description
  *markdown content*`)

	apiContent = []byte(``)
)

func TestParser_parse(t *testing.T) {
	a := assert.New(t)
	p := &parser{
		docs: make(map[string]*doc, 10),
	}

	a.NotError(p.parse(docContent))
	doc := p.docs["g1"].OpenAPI
	a.Equal(len(doc.Tags), 2)
	a.Equal(doc.Info.Title, "title")
	a.Equal(doc.Servers[0].Description, "版本号")
}
