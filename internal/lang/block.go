// SPDX-License-Identifier: MIT

package lang

import (
	"fmt"
	"strings"
	"unicode"
)

// 用于描述 block.Type 的值。
const (
	blockTypeString   int8 = iota + 1 // 字符串，将被忽略。
	blockTypeSComment                 // 单行注释
	blockTypeMComment                 // 多行注释
)

// Blocker 接口定义了解析代码块的所有操作。
// 通过 BeginFunc 查找匹配的起始位置，
// 通过 EndFunc 查找结束位置，并返回所有的块内容。
type Blocker interface {
	// 确定 l 的当前位置是否匹配 blocker 的起始位置。
	BeginFunc(l *Lexer) bool

	// 确定 l 的当前位置是否匹配 blocker 的结束位置，若匹配则返回中间的字符串。
	//
	// ok 表示是否正确匹配；
	// data 表示匹配的内容，如果不使用返回的内容，可以返回空值。
	// 比如字符串，只需要返回 true，以确保找到了结束位置，但是 data 可以直接返回 nil。
	// raw 表示匹配情况下的原始内容，data 返回的可能是经过处理，而 raw 应该是未处理的。
	//
	// 如果在到达文件末尾都没有找到结束符，则应该返回 nil, false
	EndFunc(l *Lexer) (raw, data []byte, ok bool)
}

// 定义了与语言相关的三种类型的代码块：单行注释，多行注释，字符串。
//
// block 作为 Blocker 的默认实现，能适应大部分语言的定义。
type block struct {
	// 代码块的类型，可以是字符串，单行注释或是多行注释
	Type int8

	// 块的起始字符串
	Begin string

	// 块的结束字符串，单行注释不用定义此值
	End string

	// 转义字符
	//
	// 当 Type 为 blockTypeString 时，此值表示转义字符。
	// 当 Type 为 blockTypeMComment 时，此值表示需要过滤的行首字符，比如：
	//  /**
	//   *
	//   *
	//   */
	// 以上注释，会过滤掉每一行的 * 字符。
	Escape string
}

// BeginFunc 实现 Blocker.BeginFunc
func (b *block) BeginFunc(l *Lexer) bool {
	return l.match(b.Begin)
}

// EndFunc 实现 Blocker.EndFunc
func (b *block) EndFunc(l *Lexer) (raw, data []byte, ok bool) {
	switch b.Type {
	case blockTypeString:
		return b.endString(l)
	case blockTypeMComment:
		return b.endMComments(l)
	case blockTypeSComment:
		return b.endSComments(l)
	default:
		panic(fmt.Sprintf("无效的 blockType 值：%d", b.Type))
	}
}

// 从 l 的当前位置开始往后查找，直到找到 b 中定义的 end 字符串，
// 将 l 中的指针移到该位置。
// 正常找到结束符的返回 true，否则返回 false。
//
// 第一个返回参数无用，仅是为了统一函数签名
func (b *block) endString(l *Lexer) (raw, data []byte, ok bool) {
	for {
		switch {
		case l.AtEOF():
			return nil, nil, false
		case (len(b.Escape) > 0) && l.match(b.Escape):
			l.offset++
		case l.match(b.End):
			return nil, nil, true
		default:
			l.offset++
		}
	} // end for
}

// 从 l 的当前位置往后开始查找连续的相同类型单行代码块。
func (b *block) endSComments(l *Lexer) (raw, data []byte, ok bool) {
	data = make([]byte, 0, 120)
	raw = make([]byte, 0, 120)

LOOP:
	for {
		start := l.offset // 当前行的起始位置
		for {             // 读取一行的内容
			r := l.data[l.offset]
			l.offset++
			raw = append(raw, r)

			if l.AtEOF() {
				data = append(data, l.data[start:l.offset]...)
				break LOOP
			}

			if r == '\n' {
				data = append(data, filterSymbols(l.data[start:l.offset], b.Escape)...)
				break
			}
		}

		l.skipSpace()
		if !l.match(b.Begin) { // 不是接连着的注释块了，结束当前的匹配
			break
		}
	}

	if len(data) > 0 { // 最后一个换行符返还给 Lexer
		l.offset--
	}

	return raw, data, true
}

// 从 l 的当前位置一直到定义的 b.End 之间的所有字符。
// 会对每一行应用 filterSymbols 规则。
func (b *block) endMComments(l *Lexer) (raw, data []byte, ok bool) {
	data = make([]byte, 0, 200)
	raw = make([]byte, 0, 200)
	start := l.offset

	for {
		switch {
		case l.AtEOF(): // 没有找到结束符号，直接到达文件末尾
			return nil, nil, false
		case l.match(b.End):
			if pos := l.offset - len(b.End); pos > start {
				data = append(data, filterSymbols(l.data[start:pos], b.Escape)...)
			}
			return raw, data, true
		default:
			r := l.data[l.offset]
			l.offset++
			raw = append(raw, r)
			if r == '\n' {
				data = append(data, filterSymbols(l.data[start:l.offset], b.Escape)...)
				start = l.offset
			}
		}
	} // end for
}

// 行首若出现`空白字符+symbol+空白字符`的组合，则去掉 symbol 及之前的字符。
// symbol 为 charset 中的任意字符。
func filterSymbols(line []byte, charset string) []byte {
	if len(charset) == 0 {
		return line
	}

	for index, v := range line {
		if unicode.IsSpace(rune(v)) && v != '\n' { // 跳过行首的空格，但不能换行
			continue
		}

		// 不存在指定的符号，直接返回原数据
		if strings.IndexByte(charset, v) < 0 {
			return line
		}

		// 若下个字符正好是是空格
		if len(line) > index+1 && unicode.IsSpace(rune(line[index+1])) {
			if line[index+1] == '\n' {
				return []byte{'\n'}
			}
			return line[index+1:]
		}
		return line
	}

	return line
}
