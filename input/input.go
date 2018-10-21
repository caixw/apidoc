// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package input 用于处理输入的文件，从代码中提取基本的注释内容。
//
// 多行注释和单行注释在处理上会有一定区别：
//  - 单行注释，风格相同且相邻的注释会被合并成一个注释块；
//  - 单行注释，风格不相同且相邻的注释会被按注释风格多个注释块；
//  - 多行注释，即使两个注释释块相邻也会被分成两个注释块来处理。
package input

import (
	"log"
	"math"
	"sync"
	"unicode"

	"github.com/caixw/apidoc/input/encoding"
	"github.com/caixw/apidoc/internal/locale"
)

// Block 解析出来的注释块
type Block struct {
	File string
	Line int
	Data []byte
}

// Parse 分析源代码，获取相应的文档内容。
//
// 当所有的代码块已经放入 Block 之后，Block 会被关闭。
//
// errlog 表示错误内容输出通道，若不需要，则使用 nil 代替。
func Parse(errlog *log.Logger, o ...*Options) chan Block {
	data := make(chan Block, 500)

	go func() {
		wg := &sync.WaitGroup{}
		for _, opt := range o {
			parse(data, errlog, wg, opt)
		}
		wg.Wait()

		close(data)
	}()

	return data
}

func parse(data chan Block, errlog *log.Logger, wg *sync.WaitGroup, o *Options) {
	for _, path := range o.paths {
		wg.Add(1)
		go func(path string) {
			parseFile(data, errlog, path, o)
			wg.Done()
		}(path)
	}
}

// 分析 path 指向的文件，并将内容写入到 docs 中。
//
// NOTE: parseFile 内部不能有 go 协程处理代码。
func parseFile(channel chan Block, errlog *log.Logger, path string, o *Options) {
	data, err := encoding.Transform(path, o.Encoding)
	if err != nil {
		if errlog != nil {
			errlog.Println(err)
		}
		return
	}

	l := &lexer{data: data, blocks: o.blocks}
	var block blocker

	for {
		if l.atEOF() {
			return
		}

		if block == nil {
			block = l.block()
			if block == nil { // 没有匹配的 block 了
				return
			}
		}

		ln := l.lineNumber() + 1 // 记录当前的行号，1 表示从 1 开始记数
		lines, ok := block.EndFunc(l)
		if !ok {
			errlog.Println(locale.Sprintf(locale.ErrNotFoundEndFlag))
			return // 没有找到结束标签，那肯定是到文件尾了，可以直接返回。
		}

		block = nil

		channel <- Block{
			File: path,
			Line: ln,
			Data: mergeLines(lines),
		}
	} // end for
}

// Encodings 返回支持的编码方式
func Encodings() []string {
	return encoding.Encodings()
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
