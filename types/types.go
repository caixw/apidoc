// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package types 一些公用类型的定义
package types

import "github.com/caixw/apidoc/types/openapi"

// Sanitizer 配置项的检测接口
type Sanitizer = openapi.Sanitizer

// OptionsError 提供对配置项错误的描述
type OptionsError = openapi.Error

// API 文档内容
type API struct {
	API         string              // @api 后面的内容，包含了 method, url 和 summary
	Group       string              `yaml:"group,omitempty"`
	Tags        []string            `yaml:"tags,omitempty"`
	Description openapi.Description `yaml:"description,omitempty"`
	Deprecated  bool                `yaml:"deprecated,omitempty"`
	OperationID string              `yaml:"operationId,omitempty" `
	Queries     []string            `yaml:"queries,omitempty"`
	Params      []string            `yaml:"params,omitempty"`
	Headers     []string            `yaml:"header,omitempty"`
	Request     *Request            `yaml:"request,omitempty"` // GET 此值可能为空
	Responses   []*Response         `yaml:"responses"`
}

// Request 表示请求内容
type Request struct {
	Schema   *openapi.Schema                 `yaml:"schema"`
	Examples map[string]openapi.ExampleValue `yaml:"examples,omitempty"`
}

// Response 表示返回的内容
type Response struct {
	Schema   *openapi.Schema                 `yaml:"schema"`
	Examples map[string]openapi.ExampleValue `yaml:"examples,omitempty"`
}
