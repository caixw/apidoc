// SPDX-License-Identifier: MIT

package doc

import (
	"sort"

	"github.com/caixw/apidoc/v5/internal/locale"
	"github.com/caixw/apidoc/v5/message"
)

// Sanitize 检测内容是否合法
func (doc *Doc) Sanitize() error {
	// Tag.Name 查重
	sort.SliceStable(doc.Tags, func(i, j int) bool {
		return doc.Tags[i].Name > doc.Tags[j].Name
	})
	for i := 1; i < len(doc.Tags); i++ {
		if doc.Tags[i].Name == doc.Tags[i-1].Name {
			return message.NewError(doc.file, "tag.name", doc.line, locale.ErrDuplicateValue)
		}
	}

	// TODO

	// Server.Name 查重
	sort.SliceStable(doc.Servers, func(i, j int) bool {
		return doc.Servers[i].Name > doc.Servers[j].Name
	})
	for i := 1; i < len(doc.Servers); i++ {
		if doc.Servers[i].Name == doc.Servers[i-1].Name {
			return message.NewError(doc.file, "server.name", doc.line, locale.ErrDuplicateValue)
		}
	}

	// Server.URL 查重
	sort.SliceStable(doc.Servers, func(i, j int) bool {
		return doc.Servers[i].URL > doc.Servers[j].URL
	})
	for i := 1; i < len(doc.Servers); i++ {
		if doc.Servers[i].URL == doc.Servers[i-1].URL {
			return message.NewError(doc.file, "server.url", doc.line, locale.ErrDuplicateValue)
		}
	}

	// 查看 API 中的标签是否都存在
	for _, api := range doc.Apis {
		if err := api.sanitize(); err != nil {
			return err
		}
	}

	return nil
}
