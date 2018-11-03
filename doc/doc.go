// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package doc 表示最终解析出来的文档结果。
package doc

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/issue9/is"

	"github.com/caixw/apidoc/doc/lexer"
	"github.com/caixw/apidoc/errors"
	i "github.com/caixw/apidoc/internal/input"
	"github.com/caixw/apidoc/internal/vars"
	"github.com/caixw/apidoc/options"
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

	responses

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
//
// 当所有的代码块已经放入 Block 之后，Block 会被关闭。
//
// 所有与解析有关的错误均通过 h 输出。而其它错误，比如参数问题等，通过返回参数返回。
func Parse(ctx context.Context, h *errors.Handler, input ...*options.Input) (*Doc, error) {
	start := time.Now()
	block, err := i.Parse(ctx, h, input...)
	if err != nil {
		return nil, err
	}

	doc := &Doc{
		APIDoc: vars.Version(),
	}

	wg := sync.WaitGroup{}

LOOP:
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case blk, ok := <-block:
			if !ok {
				break LOOP
			}

			wg.Add(1)
			go func(b i.Block) {
				doc.parseBlock(b, h)
				wg.Done()
			}(blk)
		}
	}

	wg.Wait()

	doc.check(h)

	doc.Elapsed = time.Now().Sub(start)

	return doc, nil
}

func (doc *Doc) check(h *errors.Handler) {
	for _, api := range doc.Apis {
		if err := api.check(); err != nil {
			h.SyntaxError(err)
		}

		for _, tag := range api.Tags {
			if !doc.tagExists(tag) {
				h.SyntaxError(api.errInvalidFormat("@apiTags")) // TODO 专门的错误信息
			}
		}

		for _, srv := range api.Servers {
			if !doc.serverExists(srv) {
				h.SyntaxError(api.errInvalidFormat("@apiServers")) // TODO 专门的错误信息
			}
		}
	}
}

func (doc *Doc) parseBlock(block i.Block, h *errors.Handler) {
	l := lexer.New(block, h)

	tag := l.Tag()
	l.Backup(tag)

	switch strings.ToLower(tag.Name) {
	case "@api":
		doc.parseAPI(l)
	case "@apidoc":
		doc.parseAPIDoc(l)
	}
}

type apiDocParser func(*Doc, *lexer.Lexer, *lexer.Tag)

// 定义了 @apidoc 子标签及其本身的的解析函数列表。
var apiDocParsers = map[string]apiDocParser{
	"@apidoc":      (*Doc).parseapidoc,
	"@apitag":      (*Doc).parseTag,
	"@apilicense":  (*Doc).parseLicense,
	"@apiserver":   (*Doc).parseServer,
	"@apicontact":  (*Doc).parseContact,
	"@apiversion":  (*Doc).parseVersion,
	"@apicontent":  (*Doc).parseContent,
	"@apiresponse": (*Doc).parseResponse,
}

func (doc *Doc) parseAPIDoc(l *lexer.Lexer) {
	for tag := l.Tag(); tag != nil; tag = l.Tag() {
		fn, found := apiDocParsers[strings.ToLower(tag.Name)]
		if !found {
			tag.ErrInvalidTag()
			return
		}

		fn(doc, l, tag)
	}
}

// 解析 @apidoc 标签，格式如下：
//  @apidoc title of document
func (doc *Doc) parseapidoc(l *lexer.Lexer, tag *lexer.Tag) {
	if doc.Title != "" {
		tag.ErrDuplicateTag()
		return
	}

	if len(tag.Data) == 0 {
		tag.ErrInvalidFormat()
		return
	}

	doc.Title = string(tag.Data)
}

// 解析 @apiContent 标签，格式如下：
//  @apicontent xxxx
func (doc *Doc) parseContent(l *lexer.Lexer, tag *lexer.Tag) {
	if doc.Content != "" {
		tag.ErrDuplicateTag()
		return
	}

	if len(tag.Data) == 0 {
		tag.ErrInvalidFormat()
		return
	}

	doc.Content = Markdown(tag.Data)
}

// 解析 @apiVersion 标签，格式如下：
//  @apiVersion 3.2.1
func (doc *Doc) parseVersion(l *lexer.Lexer, tag *lexer.Tag) {
	if doc.Version != "" {
		tag.ErrDuplicateTag()
		return
	}

	if len(tag.Data) == 0 {
		tag.ErrInvalidFormat()
		return
	}

	doc.Version = string(tag.Data)
}

// 解析 @apiContact 标签，格式如下：
//  @apiContact name name@example.com https://example.com
//  @apiContact name https://example.com
//  @apiContact name name@example.com
func (doc *Doc) parseContact(l *lexer.Lexer, tag *lexer.Tag) {
	if doc.Contact != nil {
		tag.ErrDuplicateTag()
		return
	}

	c, ok := newContact(tag)
	if ok {
		doc.Contact = c
	}
}

func (doc *Doc) tagExists(tag string) bool {
	for _, t := range doc.Tags {
		if t.Name == tag {
			return true
		}
	}

	return false
}

// 解析 @apiTag 标签，可以是以下格式
//  @apiTag admin description
func (doc *Doc) parseTag(l *lexer.Lexer, tag *lexer.Tag) {
	data := tag.Words(2)
	if len(data) != 2 {
		tag.ErrInvalidFormat()
		return
	}

	if doc.Tags == nil {
		doc.Tags = make([]*Tag, 0, 5)
	}

	name := string(data[0])
	if doc.tagExists(name) {
		tag.ErrDuplicateTag()
		return
	}

	doc.Tags = append(doc.Tags, &Tag{
		Name:        string(name),
		Description: Markdown(data[1]),
	})
}

func (doc *Doc) serverExists(name string) bool {
	for _, srv := range doc.Servers {
		if srv.Name == name {
			return true
		}
	}

	return false
}

// 解析 @apiServer 标签，可以是以下格式
//  @apiserver admin https://api1.example.com description
//  @apiserver admin https://api1.example.com
func (doc *Doc) parseServer(l *lexer.Lexer, tag *lexer.Tag) {
	data := tag.Words(3)
	if len(data) < 2 { // 描述为可选字段
		tag.ErrInvalidFormat()
		return
	}
	if !is.URL(data[1]) {
		tag.ErrInvalidFormat()
		return
	}

	if doc.Servers == nil {
		doc.Servers = make([]*Server, 0, 5)
	}

	name := string(data[0])
	if doc.serverExists(name) {
		tag.ErrDuplicateTag()
		return
	}

	srv := &Server{
		Name: name,
		URL:  string(data[1]),
	}
	if len(data) == 3 {
		srv.Description = Markdown(data[2])
	}
	doc.Servers = append(doc.Servers, srv)
}

// 解析版本信息，格式如下：
//  @apilicense MIT https://opensources.org/licenses/MIT
func (doc *Doc) parseLicense(l *lexer.Lexer, tag *lexer.Tag) {
	if doc.License != nil {
		tag.ErrDuplicateTag()
		return
	}

	data := tag.Words(2)
	if len(data) != 2 {
		tag.ErrInvalidFormat()
		return
	}

	if !is.URL(data[1]) {
		tag.ErrInvalidFormat()
		return
	}

	doc.License = &Link{
		Text: string(data[0]),
		URL:  string(data[1]),
	}
}

// 解析联系人标签内容，格式可以是：
//  @apicontact name xx@example.com https://example.com
//  @apicontact name https://example.com xx@example.com
//  @apicontact name xx@example.com
//  @apicontact name https://example.com
func newContact(tag *lexer.Tag) (c *Contact, ok bool) {
	data := tag.Words(3)

	if len(data) < 2 {
		tag.ErrInvalidFormat()
		return nil, false
	}

	contact := &Contact{Name: string(data[0])}
	v := string(data[1])
	switch checkContactType(v) {
	case 1:
		contact.URL = v
	case 2:
		contact.Email = v
	default:
		tag.ErrInvalidFormat()
		return nil, false
	}

	if len(data) == 3 {
		v := string(data[2])
		switch checkContactType(v) {
		case 1:
			contact.URL = v
		case 2:
			contact.Email = v
		default:
			tag.ErrInvalidFormat()
			return nil, false
		}
	}

	return contact, true
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
