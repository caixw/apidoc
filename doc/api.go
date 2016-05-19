// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package doc

import (
	"strings"
	"unicode"

	"github.com/caixw/apidoc/lexer"
)

// 扫描一个代码块。
//
// block 该代码块的内容；
func (doc *Doc) Scan(block string) *lexer.SyntaxError {
	l := lexer.New([]rune(block))
	return doc.scan(l)
}

// 扫描文档，生成一个Doc实例。
// 若代码块没有api文档定义，则会返回空值。
func (doc *Doc) scan(l *lexer.Lexer) *lexer.SyntaxError {
	var err *lexer.SyntaxError
	api := &API{}

LOOP:
	for {
		switch {
		case l.Match("@apiGroup "):
			api.Group = string(l.ReadLine())
			if len(api.Group) == 0 {
				return l.SyntaxError("@apiGroup未指定名称")
			}
		case l.Match("@apiQuery "):
			if api.Queries == nil {
				api.Queries = make([]*Param, 0, 1)
			}
			err = scanApiQuery(l, api)
		case l.Match("@apiParam "):
			if api.Params == nil {
				api.Params = make([]*Param, 0, 1)
			}
			p, err := scanApiParam(l)
			if err != nil {
				return err
			}
			api.Params = append(api.Params, p)
		case l.Match("@apiRequest "):
			err = scanApiRequest(l, api)
		case l.Match("@apiError "):
			resp, err := scanResponse(l)
			if err != nil {
				break
			}
			api.Error = resp
		case l.Match("@apiSuccess "):
			resp, err := scanResponse(l)
			if err != nil {
				break
			}
			api.Success = resp
		case l.Match("@api "):
			err = scanApi(l, api)
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

	// Doc的必要数据没有被初始化，说明这段代码不是api文档格式。
	if len(api.URL) == 0 || len(api.Method) == 0 {
		return nil, nil
	}

	doc.mux.Lock()
	doc.Apis = append(doc.Apis, api)
	doc.mux.Unlock()

	return nil
}

// @apiQuery size int xxxxx
func scanApiQuery(l *lexer.Lexer, api *API) *lexer.SyntaxError {
	p, err := scanApiParam(l)
	if err != nil {
		return err
	}

	api.Queries = append(api.Queries, p)
	return nil
}

func (l *lexer) scanApiRequest(d *API) error {
	r := &Request{
		Type:     l.read("\n"),
		Headers:  map[string]string{},
		Params:   []*Param{},
		Examples: []*Example{},
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
			p, err := scanApiParam(l)
			if err != nil {
				return err
			}
			r.Params = append(r.Params, p)
		case l.match("@apiExample "):
			e, err := scanApiExample(l)
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

func scanResponse(l *lexer.Lexer) (*Response, error) {
	resp := &Response{
		Headers:  map[string]string{},
		Params:   []*Param{},
		Examples: []*Example{},
	}

	words, err := l.readN(2, "\n")
	if err != nil {
		return nil, err
	}
	resp.Code = words[0]
	resp.Summary = words[1]

LOOP:
	for {
		switch {
		case l.match("@apiHeader "):
			words, err := l.readN(2, "\n")
			if err != nil {
				return nil, err
			}
			resp.Headers[words[0]] = words[1]
		case l.match("@apiParam "):
			p, err := scanApiParam(l)
			if err != nil {
				return nil, err
			}
			resp.Params = append(resp.Params, p)
		case l.match("@apiExample "):
			e, err := scanApiExample(l)
			if err != nil {
				return nil, err
			}
			resp.Examples = append(resp.Examples, e)
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

	return resp, nil
}

func scanApiExample(l *lexer.Lexer) (*Example, error) {
	words, err := l.readN(2, "@api")
	if err != nil {
		return nil, err
	}

	return &Example{
		Type: words[0],
		Code: words[1],
	}, nil
}

func scanApiParam(l *lexer.Lexer) (*Param, error) {
	p := &Param{}

	p.Name = l.ReadWord()
	p.Type = l.ReadWord()
	p.Summary = l.ReadLine()
	if len(p.Name) == 0 || len(p.Type) == 0 || len(p.Summary) == 0 {
		return nil, l.SyntaxError("缺少必要的参数")
	}
	return p, nil
}

// 若存在description参数，会原样输出，不会像其它一样去掉前导空格。
// 扫描以下格式内容：
//  @api get /test.com/api/user.json api summary
//  api description
//  api description
func scanApi(l *lexer.Lexer, api *API) error {
	words, err := l.readN(3, "\n")
	if err != nil {
		return err
	}
	api.Method = words[0]
	api.URL = words[1]
	api.Summary = words[2]
	api.Description = l.read("@api")

	return nil
}

// 将一组字符按空格进得分组，最多分 n 组
func splitN(data []rune, n int) ([]string, *lexer.SyntaxError) {
	ret := make([]string, 0, n)
	var start, end int

}

// 读取从当前位置到 delimiter 之间的所有内容，并按空格分成 n 个数组。
func (l *Lexer) ReadN(n int, delimiter string) ([]string, error) {
	ret := make([]string, 0, n)
	size := 0
	rs := []rune{}

	for {
		if l.pos >= len(l.data) || l.Match(delimiter) { // EOF或是到了下个标签处
			if delimiter != "\n" {
				l.Backup() // 若是eof，backup不会发生任何操作
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
		return nil, l.SyntaxError("未指定足够的参数")
	}
	return ret, nil
}
