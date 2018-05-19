// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package parser 解析被 input 分离出来的自定义代码块到 openapi 格式。
package parser

import (
	"bytes"
	"log"
	"sync"

	"github.com/caixw/apidoc/locale"

	yaml "gopkg.in/yaml.v2"

	"github.com/caixw/apidoc/input"
	"github.com/caixw/apidoc/openapi"
	"github.com/caixw/apidoc/vars"
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

// 添加一上新的 API 文档
func (docs *parser) newAPI(api *api) error {
	return docs.getDoc(api.Group).parseAPI(api)
}

// 添加新的文档信息
func (docs *parser) newInfo(info *Info) error {
	return docs.getDoc(info.Group).parseInfo(info)
}

// 获取指定组名的文档，group 为空，则会采用默认值组名。
func (docs *parser) getDoc(group string) *doc {
	if group == "" {
		group = vars.DefaultGroupName
	}

	docs.locker.Lock()
	defer docs.locker.Unlock()
	d, found := docs.docs[group]

	if !found {
		d = &doc{
			OpenAPI: &openapi.OpenAPI{},
		}
		docs.docs[group] = d
	}

	return d
}

// Parse 获取文档内容
func Parse(err *log.Logger, o ...*input.Options) (map[string]*openapi.OpenAPI, error) {
	p := &parser{
		docs: make(map[string]*doc, 10),
	}

	c := input.Parse(o...)

	wg := sync.WaitGroup{}
	for block := range c {
		wg.Add(1)
		go func(b input.Block) {
			defer wg.Done()
			if err := p.parse(b.Data); err != nil {
				log.Println(locale.Sprintf(locale.SyntaxError, b.File, b.Line, err.Error()))
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

func (docs *parser) parse(data []byte) error {
	data = bytes.TrimLeft(data, " ")

	if bytes.HasPrefix([]byte(vars.API), data) {
		index := bytes.IndexByte(data, '\n')
		line := data[:index]
		data = data[index+1:]
		a := &api{}
		if err := yaml.Unmarshal(data, a); err != nil {
			return err
		}

		a.API = string(line)
		return docs.newAPI(a)
	}

	if bytes.HasPrefix([]byte(vars.APIDoc), data) {
		index := bytes.IndexByte(data, '\n')
		line := data[:index]
		data = data[index+1:]
		info := &Info{}
		if err := yaml.Unmarshal(data, info); err != nil {
			return err
		}

		info.Title = string(line)
		return docs.newInfo(info)
	}

	return nil
}
