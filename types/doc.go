// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package types

import (
	"bytes"
	"sync"

	yaml "gopkg.in/yaml.v2"

	"github.com/caixw/apidoc/input"
	"github.com/caixw/apidoc/types/openapi"
	"github.com/caixw/apidoc/vars"
)

// Doc 表示单个文档
type Doc struct {
	OpenAPI *openapi.OpenAPI
	locker  sync.Mutex
}

// Docs 表示单个完整的文档实例
type Docs struct {
	// 按 group 进行分组的文档列表，
	// 每一个文档都是一个完整的 openapi 文档。
	Docs   map[string]*Doc
	locker sync.Mutex
}

// NewAPI 添加一上新的 API 文档
func (docs *Docs) NewAPI(api *API) error {
	return docs.getDoc(api.Group).parseAPI(api)
}

// NewInfo 添加新的文档信息
func (docs *Docs) NewInfo(info *Info) error {
	return docs.getDoc(info.Group).parseInfo(info)
}

// 获取指定组名的文档，group 为空，则会采用默认值组名。
func (docs *Docs) getDoc(group string) *Doc {
	if group == "" {
		group = vars.DefaultGroupName
	}

	docs.locker.Lock()
	defer docs.locker.Unlock()
	doc, found := docs.Docs[group]

	if !found {
		doc = &Doc{
			OpenAPI: &openapi.OpenAPI{},
		}
		docs.Docs[group] = doc
	}

	return doc
}

// Parse 获取文档内容
func Parse(o ...*input.Options) (map[string]*openapi.OpenAPI, error) {
	docs := &Docs{
		Docs: make(map[string]*Doc, 10),
	}

	c := input.Parse(o...)

	for block := range c {
		if err := parseData(block.Data, docs); err != nil {
			return nil, err
		}
	}

	ret := make(map[string]*openapi.OpenAPI, len(docs.Docs))
	for name, doc := range docs.Docs {
		ret[name] = doc.OpenAPI
	}

	return ret, nil
}

func parseData(data []byte, d *Docs) error {
	data = bytes.TrimLeft(data, " ")

	if bytes.HasPrefix([]byte(vars.API), data) {
		index := bytes.IndexByte(data, '\n')
		line := data[:index]
		data = data[index+1:]
		api := &API{}
		if err := yaml.Unmarshal(data, api); err != nil {
			return err
		}

		api.API = string(line)
		return d.NewAPI(api)
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
		return d.NewInfo(info)
	}

	return nil
}
