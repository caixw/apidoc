// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package doc

// 扫描文档，生成一个 Doc 实例。
//
// 若代码块没有 api 文档定义，则会返回空值。
// block 该代码块的内容；
func (doc *Doc) Scan(data []rune) *SyntaxError {
	var err *SyntaxError

	l := newLexer(data)
	api := &API{}

LOOP:
	for {
		switch {
		case l.match("@apiGroup "):
			t := l.readTag()
			api.Group = t.readWord()
			if len(api.Group) == 0 {
				return l.syntaxError("@apiGroup 未指定名称")
			}
			if !t.atEOF() {
				l.syntaxError("@apiGroup 参数过多")
			}
		case l.match("@apiQuery "):
			if api.Queries == nil {
				api.Queries = make([]*Param, 0, 1)
			}

			p, err := l.scanAPIParam()
			if err != nil {
				return err
			}
			api.Queries = append(api.Queries, p)
		case l.match("@apiParam "):
			if api.Params == nil {
				api.Params = make([]*Param, 0, 1)
			}

			p, err := l.scanAPIParam()
			if err != nil {
				return err
			}
			api.Params = append(api.Params, p)
		case l.match("@apiRequest "):
			err = l.scanAPIRequest(api)
		case l.match("@apiError "):
			api.Error, err = l.scanResponse()
		case l.match("@apiSuccess "):
			api.Success, err = l.scanResponse()
		case l.match("@api "):
			err = l.scanAPI(api)
		default:
			if l.atEOF() {
				break LOOP
			}
			l.next() // 去掉无用的字符。
		}

		if err != nil {
			return err
		}
	} // end for

	// Doc 的必要数据没有被初始化，说明这段代码不是 api 文档格式。
	if len(api.URL) == 0 || len(api.Method) == 0 {
		return nil
	}

	doc.mux.Lock()
	doc.Apis = append(doc.Apis, api)
	doc.mux.Unlock()

	return nil
}

// @apiRequest json,xml
func (l *lexer) scanAPIRequest(api *API) *SyntaxError {
	t := l.readTag()
	r := &Request{
		Type:     t.readLine(),
		Headers:  map[string]string{},
		Params:   []*Param{},
		Examples: []*Example{},
	}
	if !t.atEOF() {
		return l.syntaxError("@apiRequest 过多的参数:" + t.readEnd())
	}

LOOP:
	for {
		switch {
		case l.match("@apiHeader "):
			t := l.readTag()
			key := t.readWord()
			val := t.readLine()
			if len(key) == 0 || len(val) == 0 {
				return l.syntaxError("@apiHeader 缺少必要的参数")
			}
			if !t.atEOF() {
				return l.syntaxError("@apiHeader 参数过多")
			}
			r.Headers[string(key)] = string(val)
		case l.match("@apiParam "):
			p, err := l.scanAPIParam()
			if err != nil {
				return err
			}
			r.Params = append(r.Params, p)
		case l.match("@apiExample "):
			e, err := l.scanAPIExample()
			if err != nil {
				return err
			}
			r.Examples = append(r.Examples, e)
		case l.match("@api"): // 其它api*，退出。
			l.backup()
			break LOOP
		default:
			if l.atEOF() {
				break LOOP
			}
			l.next() // 去掉无用的字符。

		} // end switch
	} // end for

	api.Request = r
	return nil
}

func (l *lexer) scanResponse() (*Response, *SyntaxError) {
	tag := l.readTag()
	resp := &Response{
		Code:     tag.readWord(),
		Summary:  tag.readLine(),
		Headers:  map[string]string{},
		Params:   []*Param{},
		Examples: []*Example{},
	}

	if len(resp.Code) == 0 || len(resp.Summary) == 0 {
		return nil, tag.syntaxError("@apiSuccess 或是 @apiError 缺少必要的元素")
	}
	if !tag.atEOF() {
		return nil, tag.syntaxError("@apiSuccess 或是 @apiError 参数过多")
	}

LOOP:
	for {
		switch {
		case l.match("@apiHeader "):
			t := l.readTag()
			key := t.readWord()
			val := t.readLine()
			if len(key) == 0 || len(val) == 0 {
				return nil, t.syntaxError("@apiHeader 缺少必要的参数")
			}
			if !t.atEOF() {
				return nil, t.syntaxError("@apiHeader 参数过多")
			}
			resp.Headers[key] = val
		case l.match("@apiParam "):
			p, err := l.scanAPIParam()
			if err != nil {
				return nil, err
			}
			resp.Params = append(resp.Params, p)
		case l.match("@apiExample "):
			e, err := l.scanAPIExample()
			if err != nil {
				return nil, err
			}
			resp.Examples = append(resp.Examples, e)
		case l.match("@api"): // 其它api*，退出。
			l.backup()
			break LOOP
		default:
			if l.atEOF() {
				break LOOP
			}
			l.next() // 去掉无用的字符。
		}
	}

	return resp, nil
}

func (l *lexer) scanAPIExample() (*Example, *SyntaxError) {
	tag := l.readTag()
	example := &Example{
		Type: tag.readWord(),
		Code: tag.readEnd(),
	}

	if len(example.Type) == 0 || len(example.Code) == 0 {
		return nil, l.syntaxError("@apiExample 缺少必要的参数")
	}

	return example, nil
}

func (l *lexer) scanAPIParam() (*Param, *SyntaxError) {
	p := &Param{}

	tag := l.readTag()
	p.Name = tag.readWord()
	p.Type = tag.readWord()
	p.Summary = tag.readEnd()
	if len(p.Name) == 0 || len(p.Type) == 0 || len(p.Summary) == 0 {
		return nil, l.syntaxError("@apiParam 或是 @apiQuery 缺少必要的参数")
	}
	return p, nil
}

// 若存在 description 参数，会原样输出，不会像其它一样去掉前导空格。
// 扫描以下格式内容：
//  @api get /test.com/api/user.json api summary
//  api description
//  api description
func (l *lexer) scanAPI(api *API) *SyntaxError {
	t := l.readTag()
	api.Method = t.readWord()
	api.URL = t.readWord()
	api.Summary = t.readLine()

	if len(api.Method) == 0 || len(api.URL) == 0 || len(api.Summary) == 0 {
		return l.syntaxError("@api 缺少必要的参数")
	}

	api.Description = t.readEnd()

	return nil
}
