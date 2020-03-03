// SPDX-License-Identifier: MIT

package doc

import (
	"bytes"

	xmessage "golang.org/x/text/message"

	"github.com/caixw/apidoc/v6/input"
	"github.com/caixw/apidoc/v6/message"
)

func newBlockError(b *input.Block, field string, key xmessage.Reference, v ...interface{}) error {
	return message.NewLocaleError(b.File, field, b.Line, key, v...)
}

// DeleteFile 从文档中删除与文件 file 相关的文档内容
func (doc *Doc) DeleteFile(file string) {
	for index, api := range doc.Apis {
		if api.Block.File == file {
			doc.Apis = append(doc.Apis[:index], doc.Apis[index+1:]...)
		}
	}

	if doc.Block.File == file {
		doc.Block = &input.Block{}
		doc.Mimetypes = doc.Mimetypes[:0]
		doc.Title = ""
		doc.Responses = doc.Responses[:0]
		doc.APIDoc = ""
		doc.Contact = nil
		doc.Created = ""
		doc.Description = Richtext{}
		doc.Lang = ""
		doc.License = nil
		doc.Logo = ""
		doc.Servers = doc.Servers[:0]
		doc.Tags = doc.Tags[:0]
	}
}

// Parse 分析从 opt 中获取的代码块
//
// 所有与解析有关的错误均通过 h 输出。
// 如果是配置文件的错误，则通过 error 返回
func (doc *Doc) Parse(h *message.Handler, o ...*input.Options) {
	doc.parse(h, func(blocks chan input.Block) {
		input.Parse(blocks, h, o...)
	})
}

// ParseFile 分析 path 的内容，并将其中的文档解析至 doc
func (doc *Doc) ParseFile(h *message.Handler, path string, o *input.Options) {
	doc.parse(h, func(blocks chan input.Block) {
		input.ParseFile(blocks, h, path, o)
	})
}

func (doc *Doc) parse(h *message.Handler, g func(chan input.Block)) {
	done := make(chan struct{})
	blocks := make(chan input.Block, 50)

	go func() {
		for block := range blocks {
			if err := doc.ParseBlock(&block); err != nil {
				h.Error(message.Erro, err)
			}
		}
		done <- struct{}{}
	}()

	g(blocks)
	close(blocks)
	<-done
}

var (
	apidocBegin = []byte("<apidoc")
	apiBegin    = []byte("<api")
)

// ParseBlock 分析 b 的内容并填充到 doc
func (doc *Doc) ParseBlock(b *input.Block) error {
	switch {
	case bytes.HasPrefix(b.Data, apidocBegin):
		if err := doc.FromXML(b); err != nil {
			return err
		}
	case bytes.HasPrefix(b.Data, apiBegin):
		if err := doc.NewAPI(b); err != nil {
			return err
		}
	}
	return nil
}
