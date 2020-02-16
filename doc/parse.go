// SPDX-License-Identifier: MIT

package doc

import (
	"bytes"

	xmessage "golang.org/x/text/message"

	"github.com/caixw/apidoc/v6/message"
)

// Block 表示原始的注释代码块
type Block struct {
	File string
	Line int
	Data []byte
}

// NewLocaleError 生成基于 Block 定位的错误信息
func (b *Block) NewLocaleError(field string, key xmessage.Reference, v ...interface{}) error {
	return message.NewLocaleError(b.File, field, b.Line, key, v...)
}

// 从文档中删除与文件 file 相关的文档内容
func (doc *Doc) DeleteFile(file string) {
	for index, api := range doc.Apis {
		if api.Block.File == file {
			doc.Apis = append(doc.Apis[:index], doc.Apis[index+1:]...)
		}
	}

	if doc.Block.File == file {
		doc.Block = &Block{}
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

var (
	apidocBegin = []byte("<apidoc")
	apiBegin    = []byte("<api")
)

// ParseBlock 分析 b 的内容并填充到 doc
func (doc *Doc) ParseBlock(b *Block) error {
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
