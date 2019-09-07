// SPDX-License-Identifier: MIT

package doc

import (
	"bytes"
	"context"
	"encoding/xml"
	"sort"
	"sync"

	"github.com/caixw/apidoc/v5/errors"
	i "github.com/caixw/apidoc/v5/internal/input"
	"github.com/caixw/apidoc/v5/internal/locale"
	"github.com/caixw/apidoc/v5/internal/vars"
	"github.com/caixw/apidoc/v5/options"
)

// Richtext 富文本内容
type Richtext string

// Doc 文档
type Doc struct {
	XMLName struct{} `xml:"apidoc"`

	APIDoc string `xml:"-"` // 程序的版本号

	Version Version   `xml:"version,attr,omitempty"` // 文档的版本
	Title   string    `xml:"title"`
	Content Richtext  `xml:"content"`
	Contact *Contact  `xml:"contact"`
	License *Link     `xml:"license,omitempty"` // 版本信息
	Tags    []*Tag    `xml:"tag,omitempty"`     // 所有的标签
	Servers []*Server `xml:"server,omitempty"`
	Apis    []*API    `xml:"apis,omitempty"`

	// 应用于全局的变量
	//Responses []*Response `xml:"response,omitempty"`
	//Requests  []*Request  `xml:"Request,omitempty"`
	//Mimetypes string `` // 指定可用的 mimetype 类型

	references map[string]interface{}
	file       string
	line       int
}

// Tag 标签内容
type Tag struct {
	Name        string   `xml:"name,attr"`  // 字面名称，需要唯一
	Description Richtext `xml:",omitempty"` // 具体描述
	Deprecated  Version  `xml:"deprecated,attr,omitempty"`
}

// Server 服务信息
type Server struct {
	Name        string   `xml:"name,attr"` // 字面名称，需要唯一
	URL         string   `xml:"url,attr"`
	Description Richtext `xml:",omitempty"` // 具体描述
	Deprecated  Version  `xml:"deprecated,attr,omitempty"`
}

// Contact 描述联系方式
type Contact struct {
	Name  string `xml:"name,attr"`
	URL   string `xml:"url"`
	Email string `xml:"email,omitempty"`
}

// Link 表示一个链接
type Link struct {
	Text string `xml:",innerxml"`
	URL  string `xml:"url,attr"`
}

// Summary 该元素的介绍信息
type Summary struct {
	Title       string   `xml:"title,attr"`
	Description Richtext `xml:",innterxml,omitempty"`
}

// Parse 分析从 block 中获取的代码块。并填充到 Doc 中
//
// 当所有的代码块已经放入 Block 之后，Block 会被关闭。
//
// 所有与解析有关的错误均通过 h 输出。而其它错误，比如参数问题等，通过返回参数返回。
func Parse(ctx context.Context, h *errors.Handler, input ...*options.Input) (*Doc, error) {
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

	return doc, nil
}

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

func (doc *Doc) parseBlock(block i.Block, h *errors.Handler) {
	switch {
	case bytes.HasPrefix(block.Data, []byte("<apidoc ")):
		err := xml.Unmarshal(block.Data, doc)
		if serr, ok := err.(*xml.SyntaxError); ok {
			h.SyntaxError(errors.New(block.File, "", block.Line+serr.Line, serr.Msg))
		}
	case bytes.HasPrefix(block.Data, []byte("<api ")):
		api := &API{}
		err := xml.Unmarshal(block.Data, api)
		if err != nil {
			if serr, ok := err.(*xml.SyntaxError); ok {
				h.SyntaxError(errors.New(block.File, "", block.Line+serr.Line, serr.Msg))
			}
			return
		}
		api.line = block.Line
		api.file = block.File
		doc.Apis = append(doc.Apis, api)
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

func (doc *Doc) serverExists(srv string) bool {
	for _, s := range doc.Servers {
		if s.Name == srv {
			return true
		}
	}

	return false
}
