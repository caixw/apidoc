// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package lang

import (
	"log"
	"math"
	"unicode"

	"github.com/caixw/apidoc/internal/locale"
)

// Parse 分析 path 指向的文件，并返回数据
func Parse(errlog *log.Logger, data []byte, blocks []Blocker) map[int][]byte {
	l := &Lexer{data: data, Blocks: blocks}
	var block Blocker

	ret := map[int][]byte{}

	for {
		if l.AtEOF() {
			return ret
		}

		if block == nil {
			block = l.Block()
			if block == nil { // 没有匹配的 block 了
				return ret
			}
		}

		ln := l.lineNumber() + 1 // 记录当前的行号，1 表示从 1 开始记数
		lines, ok := block.EndFunc(l)
		if !ok {
			errlog.Println(locale.Sprintf(locale.ErrNotFoundEndFlag))
			return ret // 没有找到结束标签，那肯定是到文件尾了，可以直接返回。
		}

		block = nil

		ret[ln] = mergeLines(lines)
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
