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
	// TODO
	Group       string              `yaml:"group,omitempty"`
	Tags        []string            `yaml:"tags,omitempty"`
	Description openapi.Description `yaml:"description,omitempty"`
	Deprecated  bool                `yaml:"deprecated,omitempty"`
	OperationID string              `yaml:"operationId,omitempty" `
	Queries     map[string]*Query   `yaml:"queries,omitempty"`
	Params      map[string]*Param   `yaml:"params,omitempty"`
	Headers     map[string]*Header  `yaml:"header,omitempty"`
	Request     *Request            `yaml:"request,omitempty"` // GET 此值可能为空
	Responses   []*Response         `yaml:"responses"`
}

// Query 表示查询参数
type Query openapi.Parameter

// Param 表示地址中的参数信息
type Param openapi.Parameter

// Header 表示报头记录
type Header openapi.Header

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
