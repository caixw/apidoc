// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package core

import (
	"fmt"
	"strings"
	"unicode"
)

type lexer struct {
	data  []rune
	line  int    // data在file文件中行号
	file  string // 源文件名称
	pos   int    // 当前指针位置
	width int    // 最后移的字符数量
}

// line data在源文件中的起始行号
// file data所在的源文件名称
func newLexer(data []rune, line int, file string) *lexer {
	return &lexer{
		data: data,
		line: line,
		file: file,
	}
}

// 当前位置在源代码中的行号
func (l *lexer) lineNumber() int {
	count := l.line
	for i := 0; i < l.pos; i++ {
		if l.data[i] == '\n' {
			count++
		}
	}
	return count
}

// 返回一个语法错误的error接口。
func (l *lexer) syntaxError() error {
	return &SyntaxError{
		Line: l.lineNumber(),
		File: l.file,
	}
}

// 读取从当前位置到delimiter之间的所有字符,会去掉尾部空格
func (l *lexer) read(delimiter string) string {
	rs := []rune{}

	for {
		if l.pos >= len(l.data) || l.match(delimiter) { // EOF或是到了下个标签处
			if delimiter != "\n" {
				l.backup() // 若是eof，backup不会发生任何操作
			}
			break
		}
		rs = append(rs, l.data[l.pos])
		l.pos++
	} // end for

	return strings.TrimSpace(string(rs))
}

// 读取从当前位置到delimiter之间的所有内容，并按空格分成n个数组。
func (l *lexer) readN(n int, delimiter string) ([]string, error) {
	ret := make([]string, 0, n)
	size := 0
	rs := []rune{}

	for {
		if l.pos >= len(l.data) || l.match(delimiter) { // EOF或是到了下个标签处
			if delimiter != "\n" {
				l.backup() // 若是eof，backup不会发生任何操作
			}

			if len(rs) > 0 {
				// 最后一条数据，去掉尾部空格
				ret = append(ret, strings.TrimRightFunc(string(rs), unicode.IsSpace))
			}
			break
		}

		r := l.data[l.pos]
		l.pos++
		if unicode.IsSpace(r) {
			if len(rs) == 0 { // 多个连续空格
				continue
			}
			if size < n-1 {
				ret = append(ret, string(rs))
				rs = rs[:0]
				size++
				continue
			}
		}

		rs = append(rs, r)
	} // end for

	if len(ret) < n {
		return nil, l.syntaxError()
	}
	return ret, nil
}

// 判断接下去的几个字符连接起来是否正好为word，若不匹配，则不移动指针。
// 可通过lexer.backup来撤消最后一次调用。
func (l *lexer) match(word string) bool {
	if l.pos+len(word) >= len(l.data) { // 剩余字符没有word长，直接返回false
		return false
	}

	width := 0
	for _, r := range word {
		rr := l.data[l.pos]
		if rr != r {
			l.pos -= width
			return false
		}

		l.pos++
		width++
	}

	l.width = width
	return true
}

// 撤消match函数的最后次调用。指针指向执行这些函数之前的位置。
func (l *lexer) backup() {
	l.pos -= l.width
	l.width = 0
}

// 扫描文档，生成一个doc实例。
// 若代码块没有api文档定义，则会返回空值。
func (l *lexer) scan() (*doc, error) {
	d := &doc{}
	var err error

LOOP:
	for {
		switch {
		case l.match("@apiGroup "):
			err = l.scanApiGroup(d)
		case l.match("@apiQuery "):
			if d.Queries == nil {
				d.Queries = make([]*param, 0, 1)
			}
			err = l.scanApiQuery(d)
		case l.match("@apiParam "):
			if d.Params == nil {
				d.Params = make([]*param, 0, 1)
			}
			p, err := l.scanApiParam()
			if err != nil {
				return nil, err
			}
			d.Params = append(d.Params, p)
		case l.match("@apiRequest "):
			err = l.scanApiRequest(d)
		case l.match("@apiStatus "):
			if d.Status == nil {
				d.Status = make([]*status, 0, 1)
			}
			err = l.scanApiStatus(d)
		case l.match("@api "): // 放最后
			err = l.scanApi(d)
		default:
			if l.pos >= len(l.data) {
				break LOOP
			}
			l.pos++ // 去掉无用的字符。
		}

		if err != nil {
			return nil, err
		}
	} // end for

	// doc的必要数据没有被初始化，说明这段代码不是api文档格式。
	if len(d.URL) == 0 || len(d.Method) == 0 {
		return nil, fmt.Errorf("在%v:%v附近的代码并未指定@api参数")
	}

	return d, nil
}

// @apiGroup version
func (l *lexer) scanApiGroup(d *doc) error {
	d.Group = l.read("\n")
	if len(d.Group) == 0 {
		return l.syntaxError()
	}
	return nil
}

func (l *lexer) scanApiQuery(d *doc) error {
	p, err := l.scanApiParam()
	if err != nil {
		return err
	}

	d.Queries = append(d.Queries, p)
	return nil
}

func (l *lexer) scanApiRequest(d *doc) error {
	r := &request{
		Type:     l.read("\n"),
		Headers:  map[string]string{},
		Params:   []*param{},
		Examples: []*example{},
	}

LOOP:
	for {
		switch {
		case l.match("@apiHeader "):
			words, err := l.readN(2, "\n")
			if err != nil {
				return err
			}
			r.Headers[words[0]] = words[1]
		case l.match("@apiParam "):
			p, err := l.scanApiParam()
			if err != nil {
				return err
			}
			r.Params = append(r.Params, p)
		case l.match("@apiExample "):
			e, err := l.scanApiExample()
			if err != nil {
				return err
			}
			r.Examples = append(r.Examples, e)
		case l.match("@api"): // 其它api*，退出。
			l.backup()
			break LOOP
		default:
			if l.pos >= len(l.data) {
				break LOOP
			}
			l.pos++ // 去掉无用的字符。

		} // end switch
	} // end for

	d.Request = r
	return nil
}

func (l *lexer) scanApiStatus(d *doc) error {
	status := &status{
		Headers:  map[string]string{},
		Params:   []*param{},
		Examples: []*example{},
	}

	words, err := l.readN(2, "\n")
	if err != nil {
		return err
	}
	status.Code = words[0]
	status.Summary = words[1]

LOOP:
	for {
		switch {
		case l.match("@apiHeader "):
			words, err := l.readN(2, "\n")
			if err != nil {
				return err
			}
			status.Headers[words[0]] = words[1]
		case l.match("@apiParam "):
			p, err := l.scanApiParam()
			if err != nil {
				return err
			}
			status.Params = append(status.Params, p)
		case l.match("@apiExample "):
			e, err := l.scanApiExample()
			if err != nil {
				return err
			}
			status.Examples = append(status.Examples, e)
		case l.match("@api"): // 其它api*，退出。
			l.backup()
			break LOOP
		default:
			if l.pos >= len(l.data) {
				break LOOP
			}
			l.pos++ // 去掉无用的字符。
		}
	}

	d.Status = append(d.Status, status)
	return nil
}

func (l *lexer) scanApiExample() (*example, error) {
	words, err := l.readN(2, "@api")
	if err != nil {
		return nil, err
	}

	return &example{
		Type: words[0],
		Code: words[1],
	}, nil
}

func (l *lexer) scanApiParam() (*param, error) {
	words, err := l.readN(3, "\n")
	if err != nil {
		return nil, err
	}

	return &param{
		Name:        words[0],
		Type:        words[1],
		Description: words[2],
	}, nil
}

// 若存在description参数，会原样输出，不会像其它一样去掉前导空格。
// 扫描以下格式内容：
//  @api get /test.com/api/user.json api summary
//  api description
//  api description
func (l *lexer) scanApi(d *doc) error {
	words, err := l.readN(3, "\n")
	if err != nil {
		return err
	}
	d.Method = words[0]
	d.URL = words[1]
	d.Summary = words[2]
	d.Description = l.read("@api")

	return nil
}

// 扫描data，将其内容分解成doc实例，并写入到docs中
func (docs Docs) Scan(data []rune, line int, file string) error {
	l := newLexer(data, line, file)
	d, err := l.scan()
	if err != nil || d == nil {
		return err
	}

	g, found := docs[d.Group]
	if !found {
		g = make([]*doc, 0, 1)
	}

	docs[d.Group] = append(g, d)
	return nil
}
