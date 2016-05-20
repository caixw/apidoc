// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package doc

import "github.com/caixw/apidoc/lexer"

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
			if l.AtEOF() {
				break LOOP
			}
			l.Next() // 去掉无用的字符。
		}

		if err != nil {
			return err
		}
	} // end for

	// Doc的必要数据没有被初始化，说明这段代码不是api文档格式。
	if len(api.URL) == 0 || len(api.Method) == 0 {
		return l.SyntaxError("@api标签缺少必要的参数")
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

func scanApiRequest(l *lexer.Lexer, api *API) *lexer.SyntaxError {
	r := &Request{
		Type:     string(l.ReadLine()),
		Headers:  map[string]string{},
		Params:   []*Param{},
		Examples: []*Example{},
	}

LOOP:
	for {
		switch {
		case l.Match("@apiHeader "):
			key := l.ReadWord()
			val := l.ReadLine()
			if len(key) == 0 || len(val) == 0 {
				return l.SyntaxError("@apiHeader 缺少必要的参数")
			}
			r.Headers[string(key)] = string(val)
		case l.Match("@apiParam "):
			p, err := scanApiParam(l)
			if err != nil {
				return err
			}
			r.Params = append(r.Params, p)
		case l.Match("@apiExample "):
			e, err := scanApiExample(l)
			if err != nil {
				return err
			}
			r.Examples = append(r.Examples, e)
		case l.Match("@api"): // 其它api*，退出。
			l.Backup()
			break LOOP
		default:
			if l.AtEOF() {
				break LOOP
			}
			l.Next() // 去掉无用的字符。

		} // end switch
	} // end for

	api.Request = r
	return nil
}

func scanResponse(l *lexer.Lexer) (*Response, error) {
	resp := &Response{
		Headers:  map[string]string{},
		Params:   []*Param{},
		Examples: []*Example{},
	}

	resp.Code = string(l.ReadWord())
	resp.Summary = string(l.ReadLine())
	if len(resp.Code) == 0 || len(resp.Summary) == 0 {
		return nil, l.SyntaxError("缺少必要的元素")
	}

LOOP:
	for {
		switch {
		case l.Match("@apiHeader "):
			key := string(l.ReadWord())
			val := string(l.ReadLine())
			if len(key) == 0 || len(val) == 0 {
				return nil, l.SyntaxError("缺少必要的参数")
			}
			resp.Headers[key] = val
		case l.Match("@apiParam "):
			p, err := scanApiParam(l)
			if err != nil {
				return nil, err
			}
			resp.Params = append(resp.Params, p)
		case l.Match("@apiExample "):
			e, err := scanApiExample(l)
			if err != nil {
				return nil, err
			}
			resp.Examples = append(resp.Examples, e)
		case l.Match("@api"): // 其它api*，退出。
			l.Backup()
			break LOOP
		default:
			if l.AtEOF() {
				break LOOP
			}
			l.Next() // 去掉无用的字符。
		}
	}

	return resp, nil
}

func scanApiExample(l *lexer.Lexer) (*Example, *lexer.SyntaxError) {
	example := &Example{
		Type: string(l.ReadWord()),
		// TODO 多行内容
		Code: string(l.Read("@api")),
	}

	if len(example.Type) == 0 || len(example.Code) == 0 {
		return nil, l.SyntaxError("@apiExample 缺少必要的参数")
	}

	return example, nil
}

func scanApiParam(l *lexer.Lexer) (*Param, *lexer.SyntaxError) {
	p := &Param{}

	p.Name = string(l.ReadWord())
	p.Type = string(l.ReadWord())
	p.Summary = string(l.ReadLine())
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
func scanApi(l *lexer.Lexer, api *API) *lexer.SyntaxError {
	api.Method = string(l.ReadWord())
	api.URL = string(l.ReadWord())
	api.Summary = string(l.ReadLine())

	if len(api.Method) == 0 || len(api.URL) == 0 || len(api.Summary) == 0 {
		return l.SyntaxError("缺少必要的参数")
	}

	// TODO 描述内容可能为空，如何界定该内容的起止符号
	api.Description = string(l.Read("@api"))

	return nil
}
