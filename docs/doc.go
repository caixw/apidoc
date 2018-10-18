// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package docs

import (
	"bytes"

	"github.com/issue9/is"
	"github.com/issue9/version"

	"github.com/caixw/apidoc/docs/syntax"
	"github.com/caixw/apidoc/locale"
)

// @apidoc 的格式下如：
//
// @apidoc title of doc
// @apibaseurl xxxx
// @tag t1 desc
// @tag t2 desc
// @license name url
// @contact name url
// @version v1
// @external name url
// @apicontent description markdown
//
// @apiresponse 404 xxx // 全局的返回内容定义

func (docs *Docs) parseAPIDoc(l *syntax.Lexer) error {
	doc := &Doc{}
	for tag, eof := l.Tag(); !eof; tag, eof = l.Tag() {
		switch string(bytes.ToLower(tag.Name)) {
		case "@apidoc":
			if len(tag.Data) == 0 {
				return tag.Error(locale.ErrTagArgNotEnough, "@apidoc")
			}
			if doc.Title != "" {
				return tag.Error(locale.ErrDuplicateTag, "@apidoc")
			}
			doc.Title = string(tag.Data)
		case "@apitag":
			data := tag.Split(2)
			if len(data) != 2 {
				return tag.Error(locale.ErrInvalidFormat, "@apiTag")
			}
			if doc.Tags == nil {
				doc.Tags = make([]*Tag, 0, 10)
			}
			doc.Tags = append(doc.Tags, &Tag{
				Name:        string(data[0]),
				Description: Markdown(data[1]),
			})
		case "@apilicense":
			if doc.License != nil {
				return tag.Error(locale.ErrDuplicateTag, "@apiLicense")
			}

			data := tag.Split(2)
			if len(data) != 2 {
				return tag.Error(locale.ErrInvalidFormat, "@apiLicense")
			}
			if !is.URL(data[1]) {
				return tag.Error(locale.ErrInvalidFormat, "@apiLicense")
			}
			doc.License = &Link{
				Text: string(data[0]),
				URL:  string(data[1]),
			}
		case "@apicontract":
			if err := doc.parseContract(tag); err != nil {
				return err
			}
		case "@apiversion":
			if doc.Version != "" {
				return tag.Error(locale.ErrDuplicateTag, "@apiVersion")
			}
			doc.Version = string(tag.Data)

			if !version.SemVerValid(doc.Version) {
				return tag.Error(locale.ErrInvalidFormat, "@apiVersion")
			}
		case "@apibaseurl":
			doc.BaseURL = string(tag.Data)
		case "@apidescription":
			if doc.Content == "" {
				return tag.Error(locale.ErrDuplicateTag, "@apiDescription")
			}
			doc.Content = Markdown(tag.Data)
		default:
			return tag.Error(locale.ErrInvalidTag, string(tag.Name))
		}
	}

	return nil
}

func (doc *Doc) parseContract(tag *syntax.Tag) error {
	if doc.Contact != nil {
		return tag.Error(locale.ErrDuplicateTag, "@apiContract")
	}

	data := tag.Split(3)

	if len(data) < 2 || len(data) > 3 {
		return tag.Error(locale.ErrInvalidFormat, "@apiContract")
	}

	doc.Contact = &Contact{Name: string(data[0])}
	v := string(data[1])
	switch checkContractType(v) {
	case 1:
		doc.Contact.URL = v
	case 2:
		doc.Contact.Email = v
	case 3:
		return tag.Error(locale.ErrInvalidFormat, "@apiContract")
	}

	if len(data) == 3 {
		v := string(data[2])
		switch checkContractType(v) {
		case 1:
			doc.Contact.URL = v
		case 2:
			doc.Contact.Email = v
		case 3:
			return tag.Error(locale.ErrInvalidFormat, "@apiContract")
		}
	}

	return nil
}

func checkContractType(v string) int8 {
	switch {
	case is.URL(v):
		return 1
	case is.Email(v):
		return 2
	default:
		return 0
	}
}
