// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package core

import (
	"strings"
	"unicode"
)

const eof = -1

type lexer struct {
	data  []rune
	line  int    // data所在的起始行数
	file  string // 源文件名称
	pos   int    // 当前位置
	width int    // 最后移动位置的大小
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
func (l *lexer) lineNumber() (count int) {
	count = l.line
	for i := 0; i < l.pos; i++ {
		if l.data[i] == '\n' {
			count++
		}
	}
	return
}

// 返回一个语法错误的error接口。
func (l *lexer) syntaxError() error {
	return &SyntaxError{
		line: l.lineNumber(),
		file: l.file,
	}
}

// 获取下一个字符。
// 可通过lexer.backup来撤消最后一次调用。
func (l *lexer) next() rune {
	if l.pos >= len(l.data) {
		return eof
	}

	r := l.data[l.pos]
	l.pos++
	l.width = 1
	return r
}

// 读取从当前位置到换行符\n之间的内容，不包含换行符\n。
// 可通过lexer.backup来撤消最后一次调用。
func (l *lexer) nextLine() string {
	rs := []rune{} // 缓存本次操作的字符串
	width := 0     // 缓存本次操作的字符宽度，NOTE:记得在返回之前赋值给lexer.width

	for {
		if l.pos >= len(l.data) { // 提前结束
			l.width = width
			return string(rs)
		}

		r := l.data[l.pos]
		l.pos++
		width++

		if r == '\n' {
			l.width = width
			return string(rs)
		}

		rs = append(rs, r)
	} // end for
}

// 读取当前行内，当前位置到下一个空格之间的单词，
// 若当前字符为空格，则返回一个空值，且不会移动指针。
// 可通过lexer.backup来撤消最后一次调用。
func (l *lexer) nextWord() (str string, eol bool) {
	rs := []rune{}
	width := 0 // 缓存本次操作的字符宽度，NOTE:记得在返回之前赋值给lexer.width

	for {
		if l.pos >= len(l.data) {
			l.width = width
			return string(rs), true
		}

		r := l.data[l.pos]
		l.pos++
		width++

		if unicode.IsSpace(r) {
			l.pos--
			width--
			l.width = width
			return string(rs), r == '\n'
		}

		rs = append(rs, r)
	}
}

// 判断接下去的几个字符连接起来是否正好为word，若不匹配，则不移动指针。
// 可通过lexer.backup来撤消最后一次调用。
func (l *lexer) match(word string) bool {
	// 剩余字符没有word长，直接返回false
	if l.pos+len(word) >= len(l.data) {
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

// 撤消next/nextN/nextLine/nextKeyword/match函数的最后次调用。指针指向执行这些函数之前的位置。
func (l *lexer) backup() {
	l.pos -= l.width
	l.width = 0
}

// 跳过之后除换行符之外的所有空格。
func (l *lexer) skipSpace() {
	for {
		r := l.next()
		if r == eof {
			return
		}

		if !unicode.IsSpace(r) || r == '\n' {
			l.backup()
			return
		}
	} // end for
}

// 扫描文档，生成一个doc实例。
// 若代码块没有api文档定义，则会返回空值。
func (l *lexer) scan() (*doc, error) {
	d := &doc{}
	var err error

LOOP:
	for {
		switch {
		case l.match("@apiURL"):
			err = l.scanApiURL(d)
		case l.match("@apiMethods"):
			err = l.scanApiMethods(d)
		case l.match("@apiVersion"):
			err = l.scanApiVersion(d)
		case l.match("@apiGroup"):
			err = l.scanApiGroup(d)
		case l.match("@apiQuery"):
			if d.Queries == nil {
				d.Queries = make([]*param, 0, 1)
			}
			err = l.scanApiQuery(d)
		case l.match("@apiParam"):
			if d.Params == nil {
				d.Params = make([]*param, 0, 1)
			}
			p, err := l.scanApiParam()
			if err != nil {
				return nil, err
			}
			d.Params = append(d.Params, p)
		case l.match("@apiRequest"):
			err = l.scanApiRequest(d)
		case l.match("@apiStatus"):
			if d.Status == nil {
				d.Status = make([]*status, 0, 1)
			}
			err = l.scanApiStatus(d)
		case l.match("@api"): // 放最后
			err = l.scanApi(d)
		default:
			if eof == l.next() { // 去掉无用的字符。
				break LOOP
			}
		}

		if err != nil {
			return nil, err
		}
	} // end for

	// doc的必要数据没有被初始化，说明这段代码不是api文档格式。
	if len(d.URL) == 0 || len(d.Methods) == 0 {
		return nil, nil
	}

	return d, nil
}

func (l *lexer) scanApiURL(d *doc) error {
	l.skipSpace()
	str, _ := l.nextWord()
	if len(str) == 0 {
		return l.syntaxError()
	}

	d.URL = str
	return nil
}

func (l *lexer) scanApiMethods(d *doc) error {
	l.skipSpace()
	str := l.nextLine()
	if len(str) == 0 {
		return l.syntaxError()
	}

	d.Methods = str
	return nil
}

func (l *lexer) scanApiVersion(d *doc) error {
	l.skipSpace()
	str, _ := l.nextWord()
	if len(str) == 0 {
		return l.syntaxError()
	}

	d.Version = str
	return nil
}

func (l *lexer) scanApiGroup(d *doc) error {
	l.skipSpace()
	str, _ := l.nextWord()
	if len(str) == 0 {
		return l.syntaxError()
	}

	d.Group = str
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
		Type:     l.nextLine(),
		Headers:  map[string]string{},
		Params:   []*param{},
		Examples: []*example{},
	}

LOOP:
	for {
		switch {
		case l.match("@apiHeader"):
			l.skipSpace()
			key, eol := l.nextWord()
			if eol {
				return l.syntaxError()
			}

			l.skipSpace()
			val := l.nextLine()
			if len(val) == 0 {
				return l.syntaxError()
			}
			r.Headers[key] = val
		case l.match("@apiParam"):
			p, err := l.scanApiParam()
			if err != nil {
				return err
			}
			r.Params = append(r.Params, p)
		case l.match("@apiExample"):
			e, err := l.scanApiExample()
			if err != nil {
				return err
			}
			r.Examples = append(r.Examples, e)
		case l.match("@api"): // 其它api*，退出。
			l.backup()
			break LOOP
		default:
			if eof == l.next() { // 去掉无用的字符。
				break LOOP
			}
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

	var eol bool
	l.skipSpace()
	status.Code, eol = l.nextWord()
	if len(status.Code) == 0 {
		return l.syntaxError()
	}

	if !eol {
		l.skipSpace()
		status.Summary = l.nextLine()
	}

LOOP:
	for {
		switch {
		case l.match("@apiHeader"):
			l.skipSpace()
			key, eol := l.nextWord()
			if eol {
				return l.syntaxError()
			}

			l.skipSpace()
			val := l.nextLine()
			if len(val) == 0 {
				return l.syntaxError()
			}
			status.Headers[key] = val
		case l.match("@apiParam"):
			p, err := l.scanApiParam()
			if err != nil {
				return err
			}
			status.Params = append(status.Params, p)
		case l.match("@apiExample"):
			e, err := l.scanApiExample()
			if err != nil {
				return err
			}
			status.Examples = append(status.Examples, e)
		case l.match("@api"): // 其它api*，退出。
			l.backup()
			break LOOP
		default:
			if eof == l.next() { // 去掉无用的字符。
				break LOOP
			}
		}
	}

	d.Status = append(d.Status, status)
	return nil
}

func (l *lexer) scanApiExample() (*example, error) {
	e := &example{}
	rs := []rune{}

	l.skipSpace()
	e.Type, _ = l.nextWord()

	l.skipSpace()
	for {
		if l.match("@api") {
			l.backup()
			break
		}

		r := l.next()
		if r == eof {
			break
		}

		rs = append(rs, r)
	}

	e.Code = strings.TrimRightFunc(string(rs), unicode.IsSpace)
	return e, nil
}

func (l *lexer) scanApiParam() (*param, error) {
	p := &param{}
	var eol bool

	for {
		l.skipSpace()
		switch {
		case len(p.Name) == 0:
			p.Name, eol = l.nextWord()
		case len(p.Type) == 0:
			p.Type, eol = l.nextWord()
		default:
			p.Description = l.nextLine()
			eol = true
		}

		if eol {
			break
		}
	} // end for

	if len(p.Name) == 0 || len(p.Type) == 0 || len(p.Description) == 0 {
		return nil, l.syntaxError()
	}

	return p, nil
}

// 若存在description参数，会原样输出，不会像其它一样去掉前导空格。
func (l *lexer) scanApi(d *doc) error {
	l.skipSpace()
	str := l.nextLine()
	if len(str) == 0 {
		return l.syntaxError()
	}
	d.Summary = str

	rs := []rune{}
	for {
		if l.match("@api") {
			l.backup()
			break
		}

		r := l.next()
		if r == eof {
			break
		}

		rs = append(rs, r)
	}
	d.Description = string(rs)
	return nil
}

// 扫描data，将其内容分解成doc实例，并写入到tree中
func (tree *Tree) Scan(data []rune, line int, file string) error {
	l := newLexer(data, line, file)
	d, err := l.scan()
	if err != nil || d == nil {
		return err
	}

	g, found := tree.Docs[d.Group]
	if !found {
		g = make([]*doc, 0, 1)
	}

	tree.Docs[d.Group] = append(g, d)
	return nil
}
