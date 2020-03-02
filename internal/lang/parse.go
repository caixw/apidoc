// SPDX-License-Identifier: MIT

package lang

import (
	"math"
	"unicode"

	"github.com/caixw/apidoc/v6/internal/locale"
	"github.com/caixw/apidoc/v6/message"
)

// 可以作为文档的最小代码块长度
var minSize = len("<api />")

// Parse 分析 data 中的内容，并以行号作为键名，代码块作为键值返回
func Parse(file string, data []byte, blocks []Blocker, h *message.Handler) map[int][]byte {
	l := &lexer{data: data, blocks: blocks}
	var block Blocker

	ret := map[int][]byte{}

	for {
		if l.atEOF() {
			return ret
		}

		if block == nil {
			block = l.block()
			if block == nil { // 没有匹配的 block 了
				return ret
			}
		}

		ln := l.lineNumber() + 1 // 记录当前的行号，1 表示从 1 开始记数
		lines, ok := block.EndFunc(l)
		if !ok { // 没有找到结束标签，那肯定是到文件尾了，可以直接返回。
			h.Error(message.Warn, message.NewLocaleError(file, "", ln, locale.ErrNotFoundEndFlag))
			return ret
		}

		block = nil // 重置 block

		data := mergeLines(lines)
		if len(data) > minSize {
			ret[ln] = data
		}
	} // end for
}

// 合并多行为一个 []byte 结构，并去掉前导空格
func mergeLines(lines [][]byte) []byte {
	lines = trimSpaceLine(lines)

	if len(lines) == 0 {
		return nil
	}

	// 去掉第一行的所有空格
	for index, b := range lines[0] {
		if !unicode.IsSpace(rune(b)) {
			lines[0] = lines[0][index:]
			break
		}
	}

	if len(lines) == 1 {
		return lines[0]
	}

	min := math.MaxInt32
	size := 0
	for _, line := range lines[1:] {
		size += len(line)

		if isSpaceLine(line) {
			continue
		}

		for index, b := range line {
			if !unicode.IsSpace(rune(b)) {
				if min > index {
					min = index
				}
				break
			}
		}
	}

	ret := make([]byte, 0, size+len(lines[0]))
	ret = append(ret, lines[0]...)
	for _, line := range lines[1:] {
		if isSpaceLine(line) {
			ret = append(ret, line...)
		} else {
			ret = append(ret, line[min:]...)
		}
	}

	return ret
}

// 是否为空白行
func isSpaceLine(line []byte) bool {
	for _, b := range line {
		if !unicode.IsSpace(rune(b)) {
			return false
		}
	}

	return true
}

// 去掉首尾的空行
func trimSpaceLine(lines [][]byte) [][]byte {
	// 去掉开头空行
	for index, line := range lines {
		if !isSpaceLine(line) {
			lines = lines[index:]
			break
		}
	}

	// 去掉尾部的空行
	for i := len(lines) - 1; i >= 0; i-- {
		if !isSpaceLine(lines[i]) {
			lines = lines[:i+1]
			break
		}
	}

	return lines
}
