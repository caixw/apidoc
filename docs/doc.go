// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package docs

import (
	"strings"
	"sync"

	"github.com/issue9/is"
	"github.com/issue9/version"

	"github.com/caixw/apidoc/docs/lexer"
)

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

	group string
}

func (docs *Docs) parseAPIDoc(l *lexer.Lexer) error {
	doc := &Doc{}
	for tag, eof := l.Tag(); !eof; tag, eof = l.Tag() {
		switch strings.ToLower(tag.Name) {
		case "@apidoc":
			if len(tag.Data) == 0 {
				return tag.ErrInvalidFormat()
			}
			if doc.Title != "" {
				return tag.ErrDuplicateTag()
			}
			doc.Title = string(tag.Data)
		case "@apigroup":
			if doc.group != "" {
				return tag.ErrDuplicateTag()
			}
			doc.group = string(tag.Data)
		case "@apitag":
			if err := doc.parseTag(tag); err != nil {
				return err
			}
		case "@apilicense":
			if err := doc.parseLicense(tag); err != nil {
				return err
			}
		case "@apicontact":
			if err := doc.parseContact(tag); err != nil {
				return err
			}
		case "@apiversion":
			if doc.Version != "" {
				return tag.ErrDuplicateTag()
			}
			doc.Version = string(tag.Data)

			if !version.SemVerValid(doc.Version) {
				return tag.ErrInvalidFormat()
			}
		case "@apibaseurl":
			doc.BaseURL = string(tag.Data)
		case "@apicontent":
			if doc.Content == "" {
				return tag.ErrDuplicateTag()
			}
			doc.Content = Markdown(tag.Data)
		case "@apiresponse":
			if err := doc.parseResponse(l, tag); err != nil {
				return err
			}
		default:
			return tag.ErrInvalidTag()
		}
	}

	// 复制内容到 docs 中
	doc1 := docs.getDoc(doc.group)
	doc1.Title = doc.Title
	doc1.BaseURL = doc.BaseURL
	doc1.Content = doc.Content
	doc1.Contact = doc.Contact
	doc1.License = doc.License
	doc1.Version = doc.Version
	doc1.Tags = doc.Tags
	doc1.Responses = doc.Responses

	return nil
}

func (doc *Doc) parseResponse(l *lexer.Lexer, tag *lexer.Tag) error {
	if doc.Responses == nil {
		doc.Responses = make([]*Response, 10)
	}

	resp, err := newResponse(l, tag)
	if err != nil {
		return err
	}
	doc.Responses = append(doc.Responses, resp)

	return nil
}

func (doc *Doc) parseTag(tag *lexer.Tag) error {
	data := tag.Split(2)
	if len(data) != 2 {
		return tag.ErrInvalidFormat()
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

func (doc *Doc) parseLicense(tag *lexer.Tag) error {
	if doc.License != nil {
		return tag.ErrDuplicateTag()
	}

	data := tag.Split(2)
	if len(data) != 2 {
		return tag.ErrInvalidFormat()
	}
	if !is.URL(data[1]) {
		return tag.ErrInvalidFormat()
	}
	doc.License = &Link{
		Text: string(data[0]),
		URL:  string(data[1]),
	}

	return nil
}

func (doc *Doc) parseContact(tag *lexer.Tag) error {
	if doc.Contact != nil {
		return tag.ErrDuplicateTag()
	}

	data := tag.Split(3)

	if len(data) < 2 || len(data) > 3 {
		return tag.ErrInvalidFormat()
	}

	doc.Contact = &Contact{Name: string(data[0])}
	v := string(data[1])
	switch checkContractType(v) {
	case 1:
		doc.Contact.URL = v
	case 2:
		doc.Contact.Email = v
	default:
		return tag.ErrInvalidFormat()
	}

	if len(data) == 3 {
		v := string(data[2])
		switch checkContractType(v) {
		case 1:
			doc.Contact.URL = v
		case 2:
			doc.Contact.Email = v
		default:
			return tag.ErrInvalidFormat()
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
