// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package types

import (
	"sync"

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
