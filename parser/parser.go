// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package parser 解析被 input 分离出来的自定义代码块到 openapi 格式。
package parser

import (
	"bytes"
	"log"
	"sync"

	"github.com/caixw/apidoc/input"
	"github.com/caixw/apidoc/openapi"
	"github.com/caixw/apidoc/vars"
)

var (
	apiPrefix    = []byte(vars.API)
	apiDocPrefix = []byte(vars.APIDoc)
)

// 表示单个文档
type doc struct {
	OpenAPI *openapi.OpenAPI
	locker  sync.Mutex
}

type parser struct {
	// 按 group 进行分组的文档列表，
	// 每一个文档都是一个完整的 openapi 文档。
	docs   map[string]*doc
	locker sync.Mutex
}

// 获取指定组名的文档，group 为空，则会采用默认值组名。
// 不存在则创建一个新的 doc 实例
func (p *parser) getDoc(group string) *doc {
	if group == "" {
		group = vars.DefaultGroupName
	}

	p.locker.Lock()
	defer p.locker.Unlock()
	d, found := p.docs[group]

	if !found {
		d = &doc{
			OpenAPI: &openapi.OpenAPI{
				Info: &openapi.Info{},
			},
		}
		p.docs[group] = d
	}

	return d
}

// Parse 获取文档内容
func Parse(errlog *log.Logger, o ...*input.Options) (map[string]*openapi.OpenAPI, error) {

	p := &parser{
		docs: make(map[string]*doc, 10),
	}

	c := input.Parse(errlog, o...)

	wg := sync.WaitGroup{}
	for block := range c {
		wg.Add(1)
		go func(b input.Block) {
			defer wg.Done()
			if err := p.parseBlock(b); err != nil {
				errlog.Println(err)
				return
			}
		}(block)
	}
	wg.Wait()

	ret := make(map[string]*openapi.OpenAPI, len(p.docs))
	for name, doc := range p.docs {
		ret[name] = doc.OpenAPI
	}

	return ret, nil
}

func (p *parser) parseBlock(block input.Block) error {
	l := newLexer(block)

	switch {
	case bytes.HasPrefix(block.Data, []byte("@api ")):
		return p.parseAPI(l)
	case bytes.HasPrefix(block.Data, []byte("@apidoc ")):
		return p.parseAPIDoc(l)
	}

	return nil
}
