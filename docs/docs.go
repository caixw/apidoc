// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package docs 表示最终解析出来的文档结果。
package docs

import (
	"bytes"
	"log"
	"sync"

	"github.com/caixw/apidoc/docs/syntax"
	"github.com/caixw/apidoc/input"
	"github.com/caixw/apidoc/vars"
)

var (
	apiPrefix    = []byte("@api ")
	apiDocPrefix = []byte("@apidoc ")
)

// Parse 获取文档内容
func Parse(errlog *log.Logger, o ...*input.Options) (*Docs, error) {
	docs := &Docs{
		Docs:    make(map[string]*Doc, 10),
		Version: vars.Version(),
	}

	c := input.Parse(errlog, o...)

	wg := sync.WaitGroup{}
	for block := range c {
		wg.Add(1)
		go func(b input.Block) {
			defer wg.Done()
			if err := docs.parseBlock(b); err != nil {
				errlog.Println(err)
				return
			}
		}(block)
	}
	wg.Wait()

	return docs, nil
}

// 获取指定组名的文档，group 为空，则会采用默认值组名。
// 不存在则创建一个新的 doc 实例
func (docs *Docs) getDoc(group string) *Doc {
	if group == "" {
		group = vars.DefaultGroupName
	}

	docs.locker.Lock()
	defer docs.locker.Unlock()
	doc, found := docs.Docs[group]

	if !found {
		doc = &Doc{}
		docs.Docs[group] = doc
	}

	return doc
}

func (docs *Docs) parseBlock(block input.Block) error {
	l := syntax.NewLexer(block)

	switch {
	case bytes.HasPrefix(block.Data, apiPrefix):
		return docs.parseAPI(l)
	case bytes.HasPrefix(block.Data, apiDocPrefix):
		return docs.parseAPIDoc(l)
	}

	return nil
}
