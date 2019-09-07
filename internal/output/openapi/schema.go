// SPDX-License-Identifier: MIT

package openapi

import "github.com/caixw/apidoc/v5/errors"

// Schema.Type 需要的一些预定义数据类型
const (
	TypeInt      = "integer"
	TypeLong     = "long"
	TypeFloat    = "float"
	TypeDouble   = "double"
	TypeString   = "string"
	TypeBool     = "bool"
	TypePassword = "password"
)

// IsWellDataType 是否为一个正常的数据类型
func IsWellDataType(typ string) bool {
	switch typ {
	case TypeInt, TypeLong, TypeFloat, TypeDouble, TypeString, TypeBool, TypePassword:
		return true
	default:
		return false
	}
}

// Schema 定义了输出和输出的数据类型
type Schema struct {
	// TODO doc.Schema
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

// Sanitize 数据检测
func (s *Schema) Sanitize() *errors.Error {
	if s.ExternalDocs != nil {
		if err := s.ExternalDocs.Sanitize(); err != nil {
			err.Field = "externalDocs." + err.Field
			return err
		}
	}
	return nil
}
