// SPDX-License-Identifier: MIT

package openapi

import (
	"github.com/caixw/apidoc/v5/doc"
	"github.com/caixw/apidoc/v5/message"
)

// Schema.Type 需要的一些预定义数据类型
const (
	TypeInt      = "integer"
	TypeLong     = "long"
	TypeFloat    = "float"
	TypeDouble   = "double"
	TypeString   = "string"
	TypeBool     = "bool"
	TypePassword = "password"
	TypeArray    = "array"
)

func fromDocType(t doc.Type) string {
	switch string(t) {
	case doc.Number:
		return TypeInt
	case doc.String:
		return TypeString
	case doc.Bool:
		return TypeBool
	default:
		return ""
	}
}

// Schema 定义了输出和输出的数据类型
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
	AdditionalItems *Schema `json:"additionalItems,omitempty" yaml:"additionalItems,omitempty"`
	MaxItems        int     `json:"maxItems,omitempty" yaml:"maxItems,omitempty"`
	MinItems        int     `json:"minItems,omitempty" yaml:"minItems,omitempty"`
	UniqueItems     bool    `json:"uniqueItems,omitempty" yaml:"uniqueItems,omitempty"`
	Contains        *Schema `json:"contains,omitempty" yaml:"contains,omitempty"`

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

	Title         string                 `json:"title,omitempty" yaml:"title,omitempty"`
	Description   string                 `json:"description,omitempty" yaml:"description,omitempty"`
	Default       interface{}            `json:"default,omitempty" yaml:"default,omitempty"`
	ReadOnly      bool                   `json:"readOnly,omitempty" yaml:"readOnly,omitempty"`
	WriteOnly     bool                   `json:"writeOnly,omitempty" yaml:"writeOnly,omitempty"`
	Discriminator *Discriminator         `json:"discriminator,omitempty" yaml:"discriminator,omitempty"`
	XML           *XML                   `json:"xml,omitempty" yaml:"xml,omitempty"`
	ExternalDocs  *ExternalDocumentation `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
	Example       ExampleValue           `json:"example,omitempty" yaml:"example,omitempty"`
	Deprecated    bool                   `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`
}

// XML 将 Schema 转换为 XML 的相关声明
type XML struct {
	Name      string `json:"name,omitempty" yaml:"name,omitempty"`
	Namespace string `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	Prefix    string `json:"prefix,omitempty" yaml:"prefix,omitempty"`
	Attribute bool   `json:"attribute,omitempty" yaml:"attribute,omitempty"`
	Wrapped   bool   `json:"wrapped,omitempty" yaml:"wrapped,omitempty"`
}

// Discriminator Object
//
// NOTE: 暂时未用到。
type Discriminator struct {
	PropertyName string            `json:"propertyName" yaml:"propertyName"`
	Mapping      map[string]string `json:"mapping,omitempty" yaml:"mapping,omitempty"`
}

func (s *Schema) sanitize() *message.SyntaxError {
	if s.ExternalDocs != nil {
		if err := s.ExternalDocs.sanitize(); err != nil {
			err.Field = "externalDocs." + err.Field
			return err
		}
	}

	if s.Items != nil {
		if err := s.Items.sanitize(); err != nil {
			err.Field = "items." + err.Field
			return err
		}
	}

	for name, obj := range s.Properties {
		if err := obj.sanitize(); err != nil {
			err.Field = "[" + name + "]." + err.Field
			return err
		}
	}

	return nil
}

// chkArray 是否需要检测当前类型是否为数组
func newSchema(p *doc.Param, chkArray bool) *Schema {
	xml := &XML{
		Name:      p.Name,
		Namespace: p.XMLNS,
		Prefix:    p.XMLNSPrefix,
		Attribute: p.XMLAttr,
		Wrapped:   p.XMLWrapped != "",
	}

	if chkArray && p.Array {
		return &Schema{
			Type:  TypeArray,
			Items: newSchema(p, false),
			XML:   xml,
		}
	}

	s := &Schema{
		Type:        fromDocType(p.Type),
		Title:       p.Summary,
		Description: p.Description.Text,
		Default:     p.Default,
		Deprecated:  p.Deprecated != "",
		Required:    make([]string, 0, len(p.Items)),
		XML:         xml,
	}

	// enum
	if len(p.Enums) > 0 {
		s.Enum = make([]interface{}, 0, len(p.Enums))
		for _, e := range p.Enums {
			s.Enum = append(s.Enum, e.Value)
		}
	}

	// Properties / Required
	if len(p.Items) > 0 { // 如果是对象，类型改为空
		s.Type = ""
		s.Properties = make(map[string]*Schema, len(p.Items))

		for _, item := range p.Items {
			name := item.Name
			if item.Array && item.XMLWrapped != "" {
				name = item.XMLWrapped
			}

			s.Properties[name] = newSchema(item, true)
			if !item.Optional {
				s.Required = append(s.Required, item.Name)
			}
		}
	}

	return s
}

// chkArray 是否需要检测当前类型是否为数组
func newSchemaFromRequest(p *doc.Request, chkArray bool) *Schema {
	return newSchema(p.ToParam(), chkArray)
}
