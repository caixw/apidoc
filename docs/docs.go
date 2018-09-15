// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package docs 表示最终解析出来的文档结果。
package docs

// Docs 文档集合
type Docs struct {
	Docs    map[string]*Doc // 文档集，键名为分组名称
	Version string          // 当前的程序版本
}

// Doc 文档
type Doc struct {
	Title   string   `yaml:"title" json:"title"`
	BaseURL string   `yaml:"baseURL" json:"baseURL"`
	Content Markdown `yaml:"content,omitempty" json:"content,omitempty"`
	Contact *Contact `yaml:"contact,omitempty" json:"contact,omitempty"`
	License *Link    `yaml:"license,omitempty" json:"license,omitempty" ` // 版本信息
	Version string   `yaml:"version,omitempty" json:"version,omitempty"`  // 文档的版本
	Tags    []*Tag   `yaml:"tags,omitempty" json:"tags,omitempty"`        // 所有的标签
	Apis    []*API   `yaml:"apis" json:"apis"`
}

// Markdown 表示可以使用 markdown 文档
type Markdown string

// Tag 标签内容
type Tag struct {
	Name        string   `yaml:"name" json:"name"`                                   // 字面名称，需要唯一
	Description Markdown `yaml:"description,omitempty" json:"description,omitempty"` // 具体描述
}

// Contact 描述联系方式
type Contact struct {
	Name  string `yaml:"name" json:"name"`
	URL   string `yaml:"url" json:"url"`
	Email string `yaml:"email,omitempty" json:"email,omitempty"`
}

// Link 表示一个链接
type Link struct {
	Text string `yaml:"text" json:"text"`
	URL  string `yaml:"url" json:"url"`
}

// API 表示单个 API 文档
type API struct {
	Method      string      `yaml:"method" json:"method"`
	Path        string      `yaml:"path" json:"path"`
	Summary     string      `yaml:"summary" json:"summary"`
	Description Markdown    `yaml:"description,omitempty" json:"description,omitempty"`
	Tags        []*Tag      `yaml:"tags,omitempty" json:"tags,omitempty"`
	Queries     []*Param    `yaml:"queries,omitempty" json:"queries,omitempty"` // 查询参数
	Params      []*Param    `yaml:"params,omitempty" json:"params,omitempty"`   // URL 参数
	Request     *Request    `yaml:"request,omitempty" json:"request,omitempty"` // 若是 GET，则使用此描述请求的具体数据
	Responses   []*Response `yaml:"responses" json:"responses"`
	Deprecated  string      `yaml:"deprecated,omitempty" json:"deprecated,omitempty"`
}

// Request 表示用户请求所表示的数据。
type Request struct {
	Headers  map[string]string `yaml:"headers,omitempty" json:"headers,omitempty"`
	Type     *Schema           `yaml:"type" json:"type"`
	Examples []*Example        `yaml:"examples,omitempty" json:"examples,omitempty"`
}

// Response 表示一次请求或是返回的数据。
type Response struct {
	Status   string            `yaml:"status" json:"status"`
	Summary  string            `yaml:"summary,omitempty" json:"summary,omitempty"`
	Headers  map[string]string `yaml:"headers,omitempty" json:"headers,omitempty"`
	Type     *Schema           `yaml:"type" json:"type"`
	Examples []*Example        `yaml:"examples,omitempty" json:"examples,omitempty"`
}

// Param 简单参数的描述，比如查询参数等
type Param struct {
	Name     string  `yaml:"name" json:"name"`                             // 参数名称
	Type     *Schema `yaml:"type" json:"type"`                             // 类型
	Summary  string  `yaml:"summary" json:"summary"`                       // 参数介绍
	Optional bool    `yaml:"optional,omitempty" json:"optional,omitempty"` // 是否可以为空
}

// Example 示例
type Example struct {
	Summary     string   `yaml:"summary,omitempty" json:"summary,omitempty"`
	Description Markdown `yaml:"description,omitempty" json:"description,omitempty"`
	Value       string   `yaml:"value" json:"value"` // 示例内容
}
