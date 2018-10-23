// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package lang

import (
	"strings"
	"unicode"

	"github.com/caixw/apidoc/internal/locale"
)

// 用于描述 block.Type 的值。
const (
	BlockTypeNone     int8 = iota
	BlockTypeString        // 字符串，将被忽略。
	BlockTypeSComment      // 单行注释
	BlockTypeMComment      // 多行注释
)

// Blocker 接口定义了解析代码块的所有操作。
// 通过 BeginFunc 查找匹配的起始位置，
// 通过 EndFunc 查找结束位置，并返回所有的块内容。
type Blocker interface {
	// 确定 l 的当前位置是否匹配 blocker 的起始位置。
	BeginFunc(l *Lexer) bool

	// 确定 l 的当前位置是否匹配 blocker 的结束位置，若匹配则返回中间的字符串。
	// 返回内容以行为单位进行分割。
	//
	// 如果不使用返回的内容，可以返回空值。
	// 比如字符串，只需要返回 true，以确保找到了结束位置，但是内容可以直接返回 nil。
	EndFunc(l *Lexer) ([][]byte, bool)
}

// Block 定义了与语言相关的三种类型的代码块：单行注释，多行注释，字符串。
//
// Block 作为 Blocker 的默认实现，能适应大部分语言的定义。
type Block struct {
	Type   int8   // 代码块的类型，可以是字符串，单行注释或是多行注释
	Begin  string // 块的起始字符串
	End    string // 块的结束字符串，单行注释不用定义此值
	Escape string // 当 Type 为 blockTypeString 时，此值表示转义字符，Type 为其它值时，此值无意义
}

// BeginFunc 实现 Blocker.BeginFunc
func (b *Block) BeginFunc(l *Lexer) bool {
	return l.Match(b.Begin)
}

// EndFunc 实现 Blocker.EndFunc
func (b *Block) EndFunc(l *Lexer) ([][]byte, bool) {
	switch b.Type {
	case BlockTypeString:
		return b.endString(l)
	case BlockTypeMComment:
		return b.endMComments(l)
	case BlockTypeSComment:
		return b.endSComments(l)
	default:
		panic(locale.Sprintf(locale.ErrInvalidBlockType, b.Type))
	}
}

// 从 l 的当前位置开始往后查找，直到找到 b 中定义的 end 字符串，
// 将 l 中的指针移到该位置。
// 正常找到结束符的返回 true，否则返回 false。
//
// 第一个返回参数无用，仅是为了统一函数签名
func (b *Block) endString(l *Lexer) ([][]byte, bool) {
	for {
		switch {
		case l.AtEOF():
			return nil, false
		case (len(b.Escape) > 0) && l.Match(b.Escape):
			l.pos++
		case l.Match(b.End):
			return nil, true
		default:
			l.pos++
		}
	} // end for
}

// 从 l 的当前位置往后开始查找连续的相同类型单行代码块。
func (b *Block) endSComments(l *Lexer) ([][]byte, bool) {
	lines := make([][]byte, 0, 20)

LOOP:
	for {
		start := l.pos // 当前行的起始位置
		for {          // 读取一行的内容
			r := l.data[l.pos]
			l.pos++

			if l.AtEOF() {
				lines = append(lines, l.data[start:l.pos])
				break LOOP
			}

			if r == '\n' {
				lines = append(lines, filterSymbols(l.data[start:l.pos], b.Begin))
				break
			}
		}

		l.SkipSpace()
		if !l.Match(b.Begin) { // 不是接连着的注释块了，结束当前的匹配
			break
		}
	}

	if len(lines) > 0 { // 最后一个换行符返还给 Lexer
		l.pos--
	}

	return lines, true
}

// 从 l 的当前位置一直到定义的 b.End 之间的所有字符。
// 会对每一行应用 filterSymbols 规则。
func (b *Block) endMComments(l *Lexer) ([][]byte, bool) {
	lines := make([][]byte, 0, 20)
	start := l.pos

	for {
		switch {
		case l.AtEOF(): // 没有找到结束符号，直接到达文件末尾
			return nil, false
		case l.Match(b.End):
			if pos := l.pos - len(b.End); pos > start {
				lines = append(lines, filterSymbols(l.data[start:pos], b.Begin))
			}
			return lines, true
		default:
			r := l.data[l.pos]
			l.pos++
			if r == '\n' {
				lines = append(lines, filterSymbols(l.data[start:l.pos], b.Begin))
				start = l.pos
			}
		}
	} // end for
}

// 行首若出现`空白字符+symbol+空白字符`的组合，则去掉这些字符。
// symbol 为 charset 中的任意字符。
func filterSymbols(line []byte, charset string) []byte {
	for k, v := range line {
		if unicode.IsSpace(rune(v)) && v != '\n' { // 跳过行首的空格，但不能换行
			continue
		}

		// 不存在指定的符号，直接返回原数据
		if strings.IndexByte(charset, v) < 0 {
			return line
		}

		// 若下个字符正好是是空格
		if len(line) > k+1 && unicode.IsSpace(rune(line[k+1])) {
			if line[k+1] == '\n' {
				return []byte{'\n'}
			}
			return line[k+2:]
		}
		return line
	}

	return line
}
