// SPDX-License-Identifier: MIT

package ast

import (
	"sort"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/locale"
)

// Sanitize 检测内容是否合法
func (doc *APIDoc) Sanitize() error {
	// TODO

	// doc.Apis 是多线程导入的，无法保证其顺序，
	// 此处可以保证输出内容是按一定顺序排列的。
	sort.SliceStable(doc.Apis, func(i, j int) bool {
		ii := doc.Apis[i]
		jj := doc.Apis[j]

		if ii.Path.Path == jj.Path.Path {
			return ii.Method.Value.Value < jj.Method.Value.Value
		}
		return ii.Path.Path.Value.Value < jj.Path.Path.Value.Value
	})

	for _, api := range doc.Apis { // 查看 API 中的标签是否都存在
		if err := api.sanitize(); err != nil {
			return err
		}
	}

	return nil
}

func (doc *APIDoc) tagExists(tag string) bool {
	for _, s := range doc.Tags {
		if s.Name.Value.Value == tag {
			return true
		}
	}
	return false
}

func (doc *APIDoc) serverExists(srv string) bool {
	for _, s := range doc.Servers {
		if s.Name.Value.Value == srv {
			return true
		}
	}
	return false
}

// 检测和修复 api 对象，无法修复返回错误。
//
// NOTE: 需要保证 doc 已经初始化
func (api *API) sanitize() error {
	if api.doc == nil {
		panic("api.doc 未获取正确的值")
	}

	for _, tag := range api.Tags {
		if !api.doc.tagExists(tag.Content.Value) {
			loc := core.Location{URI: api.Block.Location.URI, Range: tag.Content.Range}
			return core.NewSyntaxError(loc, "", locale.ErrInvalidValue)
		}
	}

	if len(api.Servers) == 0 {
		return core.NewSyntaxError(api.Block.Location, "", locale.ErrRequired)
	}
	for _, srv := range api.Servers {
		if !api.doc.serverExists(srv.Content.Value) {
			loc := core.Location{URI: api.Block.Location.URI, Range: srv.Content.Range}
			return core.NewSyntaxError(loc, "", locale.ErrInvalidValue)
		}
	}

	return nil
}
