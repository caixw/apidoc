// SPDX-License-Identifier: MIT

package makeutil

import "bytes"

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
