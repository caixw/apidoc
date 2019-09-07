// SPDX-License-Identifier: MIT

package doc

import (
	"sort"

	"github.com/caixw/apidoc/v5/errors"
	"github.com/caixw/apidoc/v5/internal/locale"
)

// Sanitize 检测内容是否合法
func (doc *Doc) Sanitize() error {
	// Tag.Name 查重
	sort.SliceStable(doc.Tags, func(i, j int) bool {
		return doc.Tags[i].Name > doc.Tags[j].Name
	})
	for i := 1; i < len(doc.Tags); i++ {
		if doc.Tags[i].Name == doc.Tags[i-1].Name {
			return errors.New(doc.file, "tag.name", doc.line, locale.ErrDuplicateValue)
		}
	}

	// Server.Name 查重
	sort.SliceStable(doc.Servers, func(i, j int) bool {
		return doc.Servers[i].Name > doc.Servers[j].Name
	})
	for i := 1; i < len(doc.Servers); i++ {
		if doc.Servers[i].Name == doc.Servers[i-1].Name {
			return errors.New(doc.file, "server.name", doc.line, locale.ErrDuplicateValue)
		}
	}

	// Server.URL 查重
	sort.SliceStable(doc.Servers, func(i, j int) bool {
		return doc.Servers[i].URL > doc.Servers[j].URL
	})
	for i := 1; i < len(doc.Servers); i++ {
		if doc.Servers[i].URL == doc.Servers[i-1].URL {
			return errors.New(doc.file, "server.url", doc.line, locale.ErrDuplicateValue)
		}
	}

	// 查看 API 中的标签是否都存在
	for _, api := range doc.Apis {
		for _, tag := range api.Tags {
			if !doc.tagExists(tag) {
				return errors.New(api.file, "tag", api.line, locale.ErrInvalidValue)
			}
		}

		for _, srv := range api.Servers {
			if !doc.serverExists(srv) {
				return errors.New(api.file, "server", api.line, locale.ErrInvalidValue)
			}
		}
	} // end doc.Apis

	return nil
}

func (doc *Doc) tagExists(tag string) bool {
	for _, t := range doc.Tags {
		if t.Name == tag {
			return true
		}
	}

	return false
}

func (doc *Doc) serverExists(srv string) bool {
	for _, s := range doc.Servers {
		if s.Name == srv {
			return true
		}
	}

	return false
}
