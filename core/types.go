// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package core

import (
	"time"
)

type Tree struct {
	Docs map[string][]*doc `json:"docs"` // 按apiGroup分组的文档结构。
	Date time.Time         `json:"date"` // 编译时间。
}

// 表示一个api文档。
type doc struct {
	Group       string    `json:"group"`       // 所属分组
	Version     string    `json:"version"`     // 版本号
	Methods     string    `json:"methods"`     // 请求的方法，GET，POST等
	URL         string    `json:"url"`         // 请求地址
	Summary     string    `json:"summary"`     // 简要描述
	Description string    `json:"description"` // 详细描述
	Queries     []*param  `json:"queries"`     // 查询参数
	Params      []*param  `json:"params"`      // URL参数
	Request     *request  `json:"request"`     // 若是GET，则使用此描述请求的具体数据
	Status      []*status `json:"status"`      // 各种状态码下返回的数据描述
}

type request struct {
	Type     string            `json:"type"`    // 请求的类型，xml或是json
	Headers  map[string]string `json:"headers"` // 请求必须携带的头
	Params   []*param          `json:"params"`  //提交的各个字段的描述
	Examples []*example        `json:"examples"`
}

// 表示一次请求或是返回的数据。
type status struct {
	Code     string            `json:"code"`     // 状态码
	Type     string            `json:"type"`     // 提交或是返回的数据类型,xml或是json
	Summary  string            `json:"summary"`  // 该状态下的简要描述
	Headers  map[string]string `json:"headers"`  // 必须提交的头信息或是返回的头信息。
	Params   []*param          `json:"params"`   // 提交或是返回数据的各个字段描述
	Examples []*example        `json:"examples"` // 提交或是返回数据的示例，键名为数据类型，键值为示例代码
}

// 用于描述提交和返回的参数信息。
type param struct {
	Name        string `json:"name"`        // 参数名称
	Type        string `json:"type"`        // 类型
	Optional    bool   `json:"optional"`    // 可选的参数
	Description string `json:"description"` // 参数介绍
}

// 示例代码
type example struct {
	Type string `json:"type"` // 示例代码的类型，小写，xml或是json
	Code string `json:"code"` // 示例代码
}

func NewTree() *Tree {
	return &Tree{
		Docs: map[string][]*doc{},
		Date: time.Now(),
	}
}
