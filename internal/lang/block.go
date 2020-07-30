// SPDX-License-Identifier: MIT

package lang

import (
	"bytes"
	"unicode"
)

// 接口定义了解析代码块的所有操作
type blocker interface {
	// 确定 l 的当前位置是否匹配 Blocker 的起始位置。
	beginFunc(l *parser) bool

	// 确定 l 的当前位置是否匹配 Blocker 的结束位置
	//
	// data 表示匹配的内容，如果不使用返回的内容，可以返回空值。
	// 比如字符串，只需要返回 true，以确保找到了结束位置，但是 data 可以直接返回 nil。
	//
	// 如果在到达文件末尾都没有找到结束符，则应该返回 nil, false
	endFunc(l *parser) (data []byte, ok bool)
}

type (
	stringBlock struct {
		begin, end, escape string
	}

	singleComment struct {
		begin  string
		begins []byte
	}

	multipleComment struct {
		begin, end           string
		begins, ends, prefix []byte
	}
)

func newString(begin, end, escape string) blocker {
	return &stringBlock{
		begin:  begin,
		end:    end,
		escape: escape,
	}
}

func newSingleComment(begin string) blocker {
	return &singleComment{
		begin:  begin,
		begins: []byte(begin),
	}
}

func newMultipleComment(begin, end, prefix string) blocker {
	return &multipleComment{
		begin:  begin,
		end:    end,
		prefix: []byte(prefix),

		begins: []byte(begin),
		ends:   []byte(end),
	}
}

func (b *stringBlock) beginFunc(l *parser) bool {
	return l.Match(b.begin)
}

// 从 l 的当前位置开始往后查找，直到找到 b 中定义的 end 字符串，
// 将 l 中的指针移到该位置。
// 正常找到结束符的返回 true，否则返回 false。
//
// 第一个返回参数无用，仅是为了统一函数签名
func (b *stringBlock) endFunc(l *parser) (data []byte, ok bool) {
	for {
		switch {
		case l.AtEOF():
			return nil, false
		case (len(b.escape) > 0) && l.Match(b.escape):
			l.Next(1)
		case l.Match(b.end):
			return nil, true
		default:
			l.Next(1)
		}
	} // end for
}

func (b *singleComment) beginFunc(l *parser) bool {
	return l.Match(b.begin)
}

// 从 l 的当前位置往后开始查找连续的相同类型单行代码块。
func (b *singleComment) endFunc(l *parser) (data []byte, ok bool) {
	data = make([]byte, 0, 120)

	for {
		data = append(data, b.begins...)
		bs, found := l.Delim('\n', true)
		if !found { // 找不到换行符，直接填充到末尾
			data = append(data, l.All()...)
			break
		}

		data = append(data, bs...)
		data = append(data, l.Spaces('\n')...)
		if !l.Match(b.begin) { // 不是接连着的注释块了，结束当前的匹配
			break
		}
	}

	return convertSingleCommentToXML(data, b.begins), true
}

func (b *multipleComment) beginFunc(l *parser) bool {
	return l.Match(b.begin)
}

// 从 l 的当前位置一直到定义的 b.End 之间的所有字符。
// 会对每一行应用 filterSymbols 规则。
func (b *multipleComment) endFunc(l *parser) (data []byte, ok bool) {
	data, found := l.DelimString(b.end, true)
	if !found { // 没有找到结束符号，直接到达文件末尾
		return nil, false
	}

	raw := make([]byte, 0, len(b.begins)+len(data))
	raw = append(append(raw, b.begins...), data...)
	return convertMultipleCommentToXML(raw, b.begins, b.ends, b.prefix), true
}

func convertSingleCommentToXML(lines, begin []byte) []byte {
	data := make([]byte, 0, len(lines))

	newline := true
	start := -1 // 零是一个有效的数组下标
	for index, b := range lines {
		switch {
		case b == '\n':
			if start > -1 {
				for i := 0; i < index-start; i++ {
					data = append(data, ' ')
				}
				start = -1
			}
			data = append(data, b)
			newline = true
		case newline:
			switch {
			case bytes.IndexByte(begin, b) >= 0 && start == -1:
				start = index
			case unicode.IsSpace(rune(b)) && start == -1:
				data = append(data, b)
			case bytes.IndexByte(begin, b) < 0:
				if start > -1 {
					for i := 0; i < index-start; i++ { // 替换之前字符为空格
						data = append(data, ' ')
					}
					start = -1
				}
				data = append(data, b)
				newline = false
			}
		default:
			data = append(data, b)
		}
	}

	return data
}

// 转换成合法的 XML 格式
func convertMultipleCommentToXML(data, begin, end, chars []byte) []byte {
	ret := make([]byte, len(data))
	copy(ret, data)

	index := bytes.Index(data, begin)
	if index >= 0 {
		for i := range begin {
			ret[index+i] = ' '
		}
	}

	index = bytes.LastIndex(data, end)
	if index >= 0 {
		for i := range end {
			ret[index+i] = ' '
		}
	}

	return replaceSymbols(ret, chars)
}

// 替换特殊的符号为空格，使 lines 的内容为一个合法的 xml 文档
func replaceSymbols(lines, chars []byte) []byte {
	data := make([]byte, 0, len(lines))

	newline := true
	start := -1 // 零是一个有效的数组下标
	for index, b := range lines {
		switch {
		case b == '\n':
			if start > -1 {
				for i := 0; i < index-start; i++ {
					data = append(data, ' ')
				}
				start = -1
			}
			data = append(data, b)
			newline = true
		case newline:
			switch {
			case bytes.IndexByte(chars, b) >= 0 && start == -1:
				start = index
			case unicode.IsSpace(rune(b)):
				if start > -1 {
					for i := 0; i < index-start; i++ { // 替换之前字符为空格
						data = append(data, ' ')
					}
					start = -1
					newline = false
				}
				data = append(data, b)
			case bytes.IndexByte(chars, b) < 0:
				if start > -1 {
					data = append(data, lines[start:index]...)
					start = -1
				}
				data = append(data, b)
				newline = false
			}
		default:
			data = append(data, b)
		}
	}

	return data
}
