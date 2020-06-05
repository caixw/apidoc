// SPDX-License-Identifier: MIT

// Package writer 提供缓存错误的 bytes.Buffer
package writer

import "bytes"

// Writer 提供缓存错误的 bytes.Buffer
type Writer struct {
	bytes.Buffer
	Err error
}

// New 声明 Writer 实例
func New() *Writer {
	return &Writer{}
}

// WString 写入字符串
func (w *Writer) WString(str string) *Writer {
	if w.Err == nil {
		_, w.Err = w.WriteString(str)
	}
	return w
}

// WByte 写入单个字节内容
func (w *Writer) WByte(b byte) *Writer {
	return w.WBytes([]byte{b})
}

// WBytes 写入字节内容
func (w *Writer) WBytes(data []byte) *Writer {
	if w.Err == nil {
		_, w.Err = w.Write(data)
	}
	return w
}
