// SPDX-License-Identifier: MIT

package spec

import (
	"bytes"
	"encoding/xml"

	"github.com/caixw/apidoc/v7/core"
)

// DeleteFile 从文档中删除与文件 file 相关的文档内容
func (doc *APIDoc) DeleteFile(file core.URI) {
	for index, api := range doc.Apis {
		if api.Block.Location.URI == file {
			doc.Apis = append(doc.Apis[:index], doc.Apis[index+1:]...)
		}
	}

	if doc.Block.Location.URI == file {
		doc.Block = &core.Block{}
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
func (doc *APIDoc) ParseBlock(b *core.Block) error {
	index := bytes.IndexByte(b.Data, '<')
	if index < 0 {
		return nil
	}

	switch {
	case bytes.HasPrefix(b.Data[index:], apidocBegin):
		if err := doc.fromXML(b); err != nil {
			return err
		}
	case bytes.HasPrefix(b.Data[index:], apiBegin):
		if err := doc.appendAPI(b); err != nil {
			return err
		}
	}
	return nil
}

func (doc *APIDoc) fromXML(b *core.Block) error {
	doc.Block = b
	return xml.Unmarshal(b.Data, doc)
}

// appendAPI 从 b.Data 中解析新的 API 对象
func (doc *APIDoc) appendAPI(b *core.Block) error {
	api := &API{
		Block: b,
		doc:   doc,
	}
	if err := xml.Unmarshal(b.Data, api); err != nil {
		return err
	}

	doc.Apis = append(doc.Apis, api)
	return nil
}
