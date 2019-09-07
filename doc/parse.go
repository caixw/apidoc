// SPDX-License-Identifier: MIT

package doc

import (
	"sort"

	"github.com/caixw/apidoc/v5/errors"
	"github.com/caixw/apidoc/v5/internal/locale"
)

func (doc *Doc) check(h *errors.Handler) {
	// Tag.Name 查重
	sort.SliceStable(doc.Tags, func(i, j int) bool {
		return doc.Tags[i].Name > doc.Tags[j].Name
	})
	for i := 1; i < len(doc.Tags); i++ {
		if doc.Tags[i].Name == doc.Tags[i-1].Name {
			h.SyntaxError(errors.New(doc.file, "", doc.line, locale.ErrDuplicateTag))
			return
		}
	}

	// Server.Name 查重
	sort.SliceStable(doc.Servers, func(i, j int) bool {
		return doc.Servers[i].Name > doc.Servers[j].Name
	})
	for i := 1; i < len(doc.Tags); i++ {
		if doc.Servers[i].Name == doc.Servers[i-1].Name {
			h.SyntaxError(errors.New(doc.file, "", doc.line, locale.ErrDuplicateTag))
			return
		}
	}

	// Server.URL 查重
	sort.SliceStable(doc.Servers, func(i, j int) bool {
		return doc.Servers[i].URL > doc.Servers[j].URL
	})
	for i := 1; i < len(doc.Tags); i++ {
		if doc.Servers[i].URL == doc.Servers[i-1].URL {
			h.SyntaxError(errors.New(doc.file, "", doc.line, locale.ErrDuplicateTag))
			return
		}
	}

	for _, api := range doc.Apis {
		for _, tag := range api.Tags {
			if !doc.tagExists(tag) {
				h.SyntaxError(errors.New(api.file, "", api.line, locale.ErrInvalidValue))
			}
		}

		for _, srv := range api.Servers {
			if !doc.serverExists(srv) {
				h.SyntaxError(errors.New(api.file, "", api.line, locale.ErrInvalidValue))
			}
		}
	} // end doc.Apis
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
