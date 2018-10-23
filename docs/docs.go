// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package docs 表示最终解析出来的文档结果。
package docs

import (
	"log"
	"strings"
	"sync"
	"time"

	"github.com/caixw/apidoc/docs/lexer"
	"github.com/caixw/apidoc/input"
	"github.com/caixw/apidoc/internal/vars"
)

// Docs 文档集合
type Docs struct {
	Version string        `yaml:"version" json:"version"` // 当前的程序版本
	Elapsed time.Duration `yaml:"elapsed" json:"elapsed"` // 文档解析用时

	Docs   []*Doc `yaml:"docs" json:"docs"`
	locker sync.Mutex
}

// Parse 分析从 block 中获取的代码块。并填充到 Docs 中
func Parse(errlog *log.Logger, block chan input.Block) *Docs {
	docs := &Docs{
		Docs:    make([]*Doc, 0, 10),
		Version: vars.Version(),
	}

	wg := sync.WaitGroup{}
	for blk := range block {
		wg.Add(1)
		go func(b input.Block) {
			defer wg.Done()
			if err := docs.parseBlock(b); err != nil {
				errlog.Println(err)
				return
			}
		}(blk)
	}
	wg.Wait()

	return docs
}

// 获取指定组名的文档，group 为空，则返回默认组。
// 不存在则创建一个新的 doc 实例
func (docs *Docs) getDoc(group string) *Doc {
	docs.locker.Lock()
	defer docs.locker.Unlock()

	for _, doc := range docs.Docs {
		if doc.Group == group {
			return doc
		}
	}

	doc := &Doc{Group: group}
	docs.Docs = append(docs.Docs, doc)

	return doc
}

func (docs *Docs) parseBlock(block input.Block) error {
	l := lexer.New(block)

	tag, _ := l.Tag()
	l.Backup(tag)

	switch strings.ToLower(tag.Name) {
	case "@api":
		return docs.parseAPI(lexer.New(block))
	case "@apidoc":
		return docs.parseAPIDoc(lexer.New(block))
	}

	return nil
}
