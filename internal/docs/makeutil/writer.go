// SPDX-License-Identifier: MIT

package makeutil

import (
	"bytes"
	"encoding/xml"
	"io/ioutil"
	"os"

	"github.com/caixw/apidoc/v7/core"
)

const xmlFileHeader = "\n<!-- 该文件由工具自动生成，请勿手动修改！-->\n\n"

// Writer 文本输出
type Writer struct {
	bytes.Buffer
	err error
}

// NewWriter 声明 Writer 实例
func NewWriter() *Writer {
	return &Writer{}
}

// WString 写入字符串
func (w *Writer) WString(str string) *Writer {
	if w.err != nil {
		return w
	}
	_, w.err = w.WriteString(str)
	return w
}

// WBytes 写入字节内容
func (w *Writer) WBytes(data []byte) *Writer {
	if w.err != nil {
		return w
	}
	_, w.err = w.Write(data)
	return w
}

// End 结束写入过程
//
// 如果有错误，会触发 panic
func (w *Writer) End() {
	PanicError(w.err)
}

// WriteXML 将对象 v 编码成 XML 内容并写入 uri 指向的文件
func WriteXML(uri core.URI, v interface{}, indent string) error {
	data, err := xml.MarshalIndent(v, "", indent)
	if err != nil {
		return err
	}

	path, err := uri.File()
	if err != nil {
		return err
	}

	w := NewWriter()
	w.WString(xml.Header).
		WString(xmlFileHeader).
		WBytes(data).
		WString("\n") // 统一代码风格，文件末尾加一空行。

	return ioutil.WriteFile(path, w.Bytes(), os.ModePerm)
}
