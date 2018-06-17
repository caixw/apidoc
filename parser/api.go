// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package parser

import "bytes"

// @api 的格式如下：
//
// @api GET /users/{id}/logs 获取用户信息
// @group g1
// @tag t1,t2
// @version 1.0
// @deprecated desc
// @query page int default desc
// @query size int default desc
// @query state array.string [normal,lock] 状态码
// @param id int desc
// @param id int desc
//
// @request application/json {object}
// @header name desc
// @header name desc
// @param count int optional desc
// @param list array must desc
// @param list.id int optional desc
// @param list.name int must desc
// @param list.groups array.string optional desc {normal:正常,left:离职}
// @example
// {
//  count: 5,
//  list: [
//    {id:1, name: 'name1', 'group': [1,2]},
//    {id:2, name: 'name2', 'group': [1,2]}
//  ]
// }
//
// @request application/yaml {object}
//
// @response 200 application/json {array}
// @apiheader string xxx
// @param id int desc
// @param name string desc
// @param group object desc
// @param group.id int desc
//
// @response 404 application/json {object}
// @apiheader string xxx
// @param code int desc
// @param message string desc
// @param detail array.object desc
// @param detail.id string desc
// @param detail.message string desc

func (p *parser) parseAPI(l *lexer) error {
	for tag, eof := l.tag(); !eof; tag, eof = l.tag() {
		switch string(bytes.ToLower(tag.name)) {
		case "@api":
			// TODO
		}
	}
	return nil
}
