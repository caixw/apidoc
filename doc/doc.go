// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package doc 表示最终解析出来的文档结果。
package doc

import (
	"log"
	"strings"
	"sync"
	"time"

	"github.com/issue9/is"
	"github.com/issue9/version"

	"github.com/caixw/apidoc/doc/lexer"
	"github.com/caixw/apidoc/input"
	"github.com/caixw/apidoc/internal/vars"
)

// Doc 文档
type Doc struct {
	// 以下字段不对应具体的标签
	APIDoc  string        `yaml:"apidoc" json:"apidoc"`   // 当前的程序版本
	Elapsed time.Duration `yaml:"elapsed" json:"elapsed"` // 文档解析用时

	Title   string    `yaml:"title" json:"title"`
	Content Markdown  `yaml:"content,omitempty" json:"content,omitempty"`
	Contact *Contact  `yaml:"contact,omitempty" json:"contact,omitempty"`
	License *Link     `yaml:"license,omitempty" json:"license,omitempty" ` // 版本信息
	Version string    `yaml:"version,omitempty" json:"version,omitempty"`  // 文档的版本
	Tags    []*Tag    `yaml:"tags,omitempty" json:"tags,omitempty"`        // 所有的标签
	Servers []*Server `yaml:"servers,omitempty" json:"servers,omitempty"`

	// 所有接口都有可能返回的内容。
	// 比如一些错误内容的返回，可以在此处定义。
	Responses []*Response `yaml:"responses,omitempty" json:"responses,omitempty"`

	Apis   []*API `yaml:"apis" json:"apis"`
	locker sync.Mutex
}

// Markdown 表示可以使用 markdown 文档
type Markdown string

// Tag 标签内容
type Tag struct {
	Name        string   `yaml:"name" json:"name"`                                   // 字面名称，需要唯一
	Description Markdown `yaml:"description,omitempty" json:"description,omitempty"` // 具体描述
}

// Server 服务信息
type Server struct {
	Name        string   `yaml:"name" json:"name"` // 字面名称，需要唯一
	URL         string   `yaml:"url" json:"url"`
	Description Markdown `yaml:"description,omitempty" json:"description,omitempty"` // 具体描述
}

// Contact 描述联系方式
type Contact struct {
	Name  string `yaml:"name" json:"name"`
	URL   string `yaml:"url" json:"url"`
	Email string `yaml:"email,omitempty" json:"email,omitempty"`
}

// Link 表示一个链接
type Link struct {
	Text string `yaml:"text" json:"text"`
	URL  string `yaml:"url" json:"url"`
}

// Parse 分析从 block 中获取的代码块。并填充到 Doc 中
func Parse(errlog *log.Logger, block chan input.Block) *Doc {
	doc := &Doc{
		APIDoc: vars.Version(),
	}

	wg := sync.WaitGroup{}
	for blk := range block {
		wg.Add(1)
		go func(b input.Block) {
			defer wg.Done()
			if err := doc.parseBlock(b); err != nil {
				errlog.Println(err)
				return
			}
		}(blk)
	}
	wg.Wait()

	return doc
}

func (doc *Doc) parseBlock(block input.Block) error {
	l := lexer.New(block)

	tag, _ := l.Tag()
	l.Backup(tag)

	switch strings.ToLower(tag.Name) {
	case "@api":
		return doc.parseAPI(lexer.New(block))
	case "@apidoc":
		return doc.parseAPIDoc(lexer.New(block))
	}

	return nil
}

func (doc *Doc) parseAPIDoc(l *lexer.Lexer) (err error) {
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
		case "@apitag":
			if err = doc.parseTag(tag); err != nil {
				return err
			}
		case "@apilicense":
			if err = doc.parseLicense(tag); err != nil {
				return err
			}
		case "@apiserver":
			if err = doc.parseServer(tag); err != nil {
				return err
			}
		case "@apicontact":
			if doc.Contact != nil {
				return tag.ErrDuplicateTag()
			}

			if doc.Contact, err = newContact(tag); err != nil {
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
		case "@apicontent":
			if doc.Content != "" {
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

	return nil
}

// 添加一个 API 实例
func (doc *Doc) append(api *API) {
	doc.locker.Lock()
	doc.Apis = append(doc.Apis, api)
	doc.locker.Unlock()
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

// 解析 @apiTag 标签，可以是以下格式
//  @apiTag admin description
func (doc *Doc) parseTag(tag *lexer.Tag) error {
	data := tag.Words(2)
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

// 解析 @apiServer 标签，可以是以下格式
//  @apiserver admin https://api1.example.com description
//  @apiserver admin https://api1.example.com
func (doc *Doc) parseServer(tag *lexer.Tag) error {
	data := tag.Words(3)
	if len(data) < 2 { // 描述为可选字段
		return tag.ErrInvalidFormat()
	}

	if doc.Servers == nil {
		doc.Servers = make([]*Server, 0, 5)
	}

	srv := &Server{
		Name: string(data[0]),
		URL:  string(data[1]),
	}
	if len(data) == 3 {
		srv.Description = Markdown(data[2])
	}
	doc.Servers = append(doc.Servers, srv)

	return nil
}

func (doc *Doc) parseLicense(tag *lexer.Tag) (err error) {
	if doc.License != nil {
		return tag.ErrDuplicateTag()
	}

	doc.License, err = newLink(tag)
	return err
}

// 解析链接元素，格式如下：
//  @tag text https://example.com
func newLink(tag *lexer.Tag) (*Link, error) {
	data := tag.Words(2)
	if len(data) != 2 {
		return nil, tag.ErrInvalidFormat()
	}

	if !is.URL(data[1]) {
		return nil, tag.ErrInvalidFormat()
	}

	return &Link{
		Text: string(data[0]),
		URL:  string(data[1]),
	}, nil
}

// 解析联系人标签内容，格式可以是：
//  @apicontact name xx@example.com https://example.com
//  @apicontact name https://example.com xx@example.com
//  @apicontact name xx@example.com
//  @apicontact name https://example.com
func newContact(tag *lexer.Tag) (*Contact, error) {
	data := tag.Words(3)

	if len(data) < 2 {
		return nil, tag.ErrInvalidFormat()
	}

	contact := &Contact{Name: string(data[0])}
	v := string(data[1])
	switch checkContactType(v) {
	case 1:
		contact.URL = v
	case 2:
		contact.Email = v
	default:
		return nil, tag.ErrInvalidFormat()
	}

	if len(data) == 3 {
		v := string(data[2])
		switch checkContactType(v) {
		case 1:
			contact.URL = v
		case 2:
			contact.Email = v
		default:
			return nil, tag.ErrInvalidFormat()
		}
	}

	return contact, nil
}

func checkContactType(v string) int8 {
	switch {
	case is.Email(v): // Email 也属于一种 URL
		return 2
	case is.URL(v):
		return 1
	default:
		return 0
	}
}
