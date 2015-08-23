// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package core

import (
	"fmt"
	"sync"
)

// 返回data的第一个注释块，及之后字符的位置。
type ScanFunc func(data []byte) ([]rune, int)

// Docs表示一个项目的完整文档列表。
type Docs struct {
	items []*Doc
	errs  []error
	mux   sync.Mutex
}

// 表示一个api文档。
type Doc struct {
	Group       string    // 所属分组
	Method      string    // 请求的方法，GET，POST等
	URL         string    // 请求地址
	Summary     string    // 简要描述
	Description string    // 详细描述
	Queries     []*Param  // 查询参数
	Params      []*Param  // URL参数
	Request     *Request  // 若是GET，则使用此描述请求的具体数据
	Success     *Response // 成功时的响应内容
	Error       *Response // 出错时的响应内容
}

// 表示api请求数据
type Request struct {
	Type     string            // 请求的类型，xml或是json
	Headers  map[string]string // 请求必须携带的头
	Params   []*Param          //提交的各个字段的描述
	Examples []*Example
}

// 表示一次请求或是返回的数据。
type Response struct {
	Code     string            // http状态码
	Summary  string            // 该状态下的简要描述
	Headers  map[string]string // 必须提交的头信息或是返回的头信息。
	Params   []*Param          // 提交或是返回数据的各个字段描述
	Examples []*Example        // 提交或是返回数据的示例，键名为数据类型，键值为示例代码
}

// 用于描述提交和返回的参数信息。
type Param struct {
	Name    string // 参数名称
	Type    string // 类型
	Summary string // 参数介绍
}

// 示例代码
type Example struct {
	Type string // 示例代码的类型，小写，xml或是json
	Code string // 示例代码
}

type SyntaxError struct {
	Line    int
	File    string
	Message string
}

// 判断是否存在语法错误。
func (d *Docs) HasError() bool {
	return len(d.errs) > 0
}

// 返回所有的Doc实例。
func (d *Docs) Items() []*Doc {
	return d.items
}

// 获取语法错误列表。
func (d *Docs) Errors() []error {
	return d.errs
}

func (err *SyntaxError) Error() string {
	return fmt.Sprintf("%v@%v:%v", err.Message, err.File, err.Line)
}
