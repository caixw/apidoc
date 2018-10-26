// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package schema 定义 JSON Schema 的相关操作
package schema

import (
	"bytes"

	"github.com/caixw/apidoc/doc/lexer"
)

// Schema.Type 的值枚举
const (
	Null    = "null"
	Bool    = "boolean"
	Object  = "object"
	Array   = "array"
	Number  = "number"
	String  = "string"
	Integer = "integer"
)

// Schema 简化的 JSON Schema
// https://json-schema.org/latest/json-schema-validation.html
type Schema struct {
	Type string        `json:"type,omitempty" yaml:"type,omitempty"`
	Enum []interface{} `json:"enum,omitempty" yaml:"enum,omitempty"`

	// 数值验证
	MultipleOf       int  `json:"multipleOf,omitempty" yaml:"multipleOf,omitempty"`
	Maximum          int  `json:"maximum,omitempty" yaml:"maximum,omitempty"`
	ExclusiveMaximum bool `json:"exclusiveMaximum,omitempty" yaml:"exclusiveMaximum,omitempty"`
	Minimum          int  `json:"minimum,omitempty" yaml:"minimum,omitempty"`
	ExclusiveMinimum bool `json:"exclusiveMinimum,omitempty" yaml:"exclusiveMinimum,omitempty"`

	// 字符串验证
	MaxLength int    `json:"maxLength,omitempty" yaml:"maxLength,omitempty"`
	MinLength int    `json:"minLength,omitempty" yaml:"minLength,omitempty"`
	Pattern   string `json:"pattern,omitempty" yaml:"pattern,omitempty"`

	// 数组验证
	Items           *Schema `json:"items,omitempty" yaml:"items,omitempty"`
	AdditionalItems *Schema `json:"additionalItems,omitempty" ymal:"additionalItems,omitempty"`
	MaxItems        int     `json:"maxItems,omitempty" yaml:"maxItems,omitempty"`
	MinItems        int     `json:"minItems,omitempty" yaml:"minItems,omitempty"`
	UniqueItems     bool    `json:"uniqueItems,omitempty" yaml:"uniqueItems,omitempty"`
	Contains        *Schema `json:"contains,omtempty" yaml:"contains,omtempty"`

	// 对象验证
	MaxProperties        int                `json:"maxProperties,omitempty" yaml:"maxProperties,omitempty"`
	MinProperties        int                `json:"minProperties,omitempty" yaml:"minProperties,omitempty"`
	Required             []string           `json:"required,omitempty" yaml:"required,omitempty"`
	Properties           map[string]*Schema `json:"properties,omitempty" yaml:"properties,omitempty"`
	PatternProperties    map[string]*Schema `json:"patternProperties,omitempty" yaml:"patternProperties,omitempty"`
	AdditionalProperties map[string]*Schema `json:"additionalProperties,omitempty" yaml:"additionalProperties,omitempty"`
	Dependencies         map[string]*Schema `json:"dependencies,omitempty" yaml:"dependencies,omitempty"`
	PropertyNames        *Schema            `json:"propertyNames,omitempty" yaml:"propertyNames,omitempty"`

	AllOf []*Schema `json:"allOf,omitempty" yaml:"allOf,omitempty"`
	AnyOf []*Schema `json:"anyOf,omitempty" yaml:"anyOf,omitempty"`
	OneOf []*Schema `json:"oneOf,omitempty" yaml:"oneOf,omitempty"`
	Not   *Schema   `json:"not,omitempty" yaml:"not,omitempty"`

	// 可复用对象的定义
	Definitions map[string]*Schema `json:"definitions,omitempty" yaml:"definitions,omitempty"`
	Ref         string             `json:"$ref,omitempty" yaml:"$ref,omitempty"`

	Title       string      `json:"title,omitempty" yaml:"title,omitempty"`
	Description string      `json:"description,omitempty" yaml:"description,omitempty"`
	Default     interface{} `json:"default,omitempty" yaml:"default,omitempty"`
	ReadOnly    bool        `json:"readOnly,omitempty" yaml:"readOnly,omitempty"`
	WriteOnly   bool        `json:"writeOnly,omitempty" yaml:"writeOnly,omitempty"`
}

var seqaratorDot = []byte{'.'}

// Build 用于将一条语名添加到 Schema 对象，作为其字段，语句可能是以下格式：
// @param list.groups array.string [locked,deleted] desc markdown
//  * xx: xxxxx
//  * xx: xxxxx
//
//
// name 表示变量的名称。若为空，表示是顶层的对象。
// 若子元素，则需要多层嵌套，比如：
//  list.user.id
//
// typ 表示类型中的内容，比如：
//  array, object, array.string
//
// optional 表示可选参数中的描述内容，有以下三种方式：
//  - optional 表示可选，默认为零值
//  - xx 表示可选，默认值为 xx
//  - required 表示必须
func (schema *Schema) Build(tag *lexer.Tag, name, typ, optional, desc []byte) error {
	type0, type1, err := parseType(tag, typ)
	if err != nil {
		return err
	}

	var p *Schema
	var last []byte // 最后的名称
	if len(name) > 0 {
		names := bytes.Split(name, seqaratorDot)
		for _, name := range names {
			if schema.Properties == nil {
				schema.Properties = make(map[string]*Schema, 2)
			}

			ss := schema.Properties[string(name)]
			if ss == nil {
				ss = &Schema{}
				schema.Properties[string(name)] = ss
			}
			p = schema
			last = name
			schema = ss
		}
	}

	schema.Type = type0
	schema.Description = string(desc)
	if type0 == Array {
		schema.Items = &Schema{Type: type1}
	}

	opt, def, err := parseOptional(type0, type1, optional)
	if err != nil {
		return err
	}

	if !opt {
		if p != nil {
			if p.Required == nil {
				p.Required = make([]string, 0, 10)
			}
			p.Required = append(p.Required, string(last))
		}
	} else {
		schema.Default = def
	}

	schema.Enum, err = parseEnum(schema.Type, desc)
	return err
}
