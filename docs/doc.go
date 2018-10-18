// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package docs

import (
	"bytes"
	"sync"

	"github.com/issue9/is"
	"github.com/issue9/version"

	"github.com/caixw/apidoc/docs/syntax"
	"github.com/caixw/apidoc/locale"
)

// @apidoc 的格式下如：
//
// @apidoc title of doc
// @apibaseurl xxxx
// @apitag t1 desc
// @apitag t2 desc
// @apilicense name url
// @apicontact name url
// @apiversion v1
// @apicontent description markdown
//
// @apiresponse 404 xxx // 全局的返回内容定义

// Doc 文档
type Doc struct {
	Title   string   `yaml:"title" json:"title"`
	BaseURL string   `yaml:"baseURL" json:"baseURL"`
	Content Markdown `yaml:"content,omitempty" json:"content,omitempty"`
	Contact *Contact `yaml:"contact,omitempty" json:"contact,omitempty"`
	License *Link    `yaml:"license,omitempty" json:"license,omitempty" ` // 版本信息
	Version string   `yaml:"version,omitempty" json:"version,omitempty"`  // 文档的版本
	Tags    []*Tag   `yaml:"tags,omitempty" json:"tags,omitempty"`        // 所有的标签

	// 所有接口都有可能返回的内容。
	// 比如一些错误内容的返回，可以在此处定义。
	Responses []*Response `yaml:"responses,omitempty" json:"responses,omitempty"`

	Apis   []*API `yaml:"apis" json:"apis"`
	locker sync.Mutex
}

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
			if err := doc.parseTag(tag); err != nil {
				return err
			}
		case "@apilicense":
			if err := doc.parseLicense(tag); err != nil {
				return err
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
		case "@apicontent":
			if doc.Content == "" {
				return tag.Error(locale.ErrDuplicateTag, "@apiContent")
			}
			doc.Content = Markdown(tag.Data)
		case "@apiresponse":
			// TODO
		default:
			return tag.Error(locale.ErrInvalidTag, string(tag.Name))
		}
	}

	return nil
}

func (doc *Doc) parseTag(tag *syntax.Tag) error {
	data := tag.Split(2)
	if len(data) != 2 {
		return tag.Error(locale.ErrInvalidFormat, "@apiTag")
	}

	if doc.Tags == nil {
		doc.Tags = make([]*Tag, 0, 5)
	}

	doc.Tags = append(doc.Tags, &Tag{
		Name:        string(data[0]),
		Description: Markdown(data[1]),
	})

	return nil
}

func (doc *Doc) parseLicense(tag *syntax.Tag) error {
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
