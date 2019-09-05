// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package doc

import (
	"bufio"
	"bytes"
	"sort"
	"strconv"
	"unicode"

	"github.com/caixw/apidoc/errors"
	"github.com/caixw/apidoc/internal/locale"
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

	Title       string      `json:"title,omitempty" yaml:"title,omitempty"`
	Description string      `json:"description,omitempty" yaml:"description,omitempty"`
	Default     interface{} `json:"default,omitempty" yaml:"default,omitempty"`
	ReadOnly    bool        `json:"readOnly,omitempty" yaml:"readOnly,omitempty"`
	WriteOnly   bool        `json:"writeOnly,omitempty" yaml:"writeOnly,omitempty"`
}

var seqaratorDot = []byte{'.'}

// 用于将一条语名添加到 Schema 对象，作为其子元素，语句可能是以下格式：
// @param list.groups array.string [locked,deleted] desc markdown
//  - xx xxxxx
//  - xx xxxxx
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
func (schema *Schema) build(name, typ, optional, desc []byte) error {
	type0, type1, err := parseType(typ)
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

// 分析类型的内容。值可以有以下格式：
//  - type 单一类型
//  - type.subtype 集合类型，subtype 表示集全元素的类型，一般用于数组。
func parseType(typ []byte) (t1, t2 string, err error) {
	types := bytes.SplitN(typ, seqaratorDot, 2)
	if len(types) == 0 {
		return "", "", &errors.LocaleError{MessageKey: locale.ErrInvalidFormat}
	}

	type0 := string(types[0])
	if type0 != Array {
		return type0, "", nil
	}

	if len(types) == 1 {
		return "", "", &errors.LocaleError{MessageKey: locale.ErrInvalidFormat}
	}

	return type0, string(types[1]), nil
}

// 分析可选类型，格式如下
//  optional   // 默认值为零值
//  required   // 必填
//  optional.defaultvalue // 选填，默认值为 defaultvalue
func parseOptional(typ, subtype string, optional []byte) (opt bool, val interface{}, err error) {
	index := bytes.IndexByte(optional, '.')
	if index < 0 {
		switch typ {
		case Array:
			val = []interface{}{}
		case Number, Integer:
			val = 0
		case String:
			val = ""
		case Bool:
			val = false
		}
		return isOptional(optional), val, nil
	}

	opt = isOptional(optional[:index])

	optional = optional[index+1:]
	if typ == Array {
		fn := getConvertFunc(subtype)
		data := parseArray(optional)
		vals := make([]interface{}, 0, len(data))
		for _, item := range data {
			v, err := fn(string(item))
			if err != nil {
				return false, nil, err
			}
			vals = append(vals, v)
		}
		val = vals
	} else {
		fn := getConvertFunc(typ)
		if val, err = fn(string(optional)); err != nil {
			return false, nil, err
		}
	}

	return opt, val, nil
}

// 解析数组
//  "[a1,a2,a3]" ==> {"a1","a2","a3"}
//  "a1" ==> {"a1"}
func parseArray(optional []byte) [][]byte {
	optional = bytes.TrimFunc(optional, func(r rune) bool {
		return r == '[' || r == ']'
	})

	ret := make([][]byte, 0, bytes.Count(optional, []byte{','}))
LOOP:
	for {
		index := bytes.IndexByte(optional, ',')
		switch {
		case index < 0:
			ret = append(ret, bytes.TrimSpace(optional))
			break LOOP
		case index > 0:
			ret = append(ret, bytes.TrimSpace(optional[:index]))
			optional = optional[index+1:]
		case index == 0:
			ret = append(ret, []byte{})
			optional = optional[1:]
		}
	}

	return ret
}

// 分析枚举内容
//
// type 表示希望最终返回的类型。
//
// data 支持类似的的格式：
//  - s1 xxxx
//  - s2 xx
//  多行内容
//  - s3 xxxx
// 将返回 s1,s2,s3
func parseEnum(typ string, data []byte) ([]interface{}, error) {
	enum := make([]string, 0, 5)

	scanner := bufio.NewScanner(bytes.NewReader(data))
	scanner.Split(bufio.ScanLines)

LOOP:
	for scanner.Scan() {
		line := scanner.Bytes()
		line = bytes.TrimSpace(line)

		// 过滤非 - 开头的行
		//
		// bug(caixw):如果需要添加其它类型的前缀符号，
		// 请注意 internal/lang.filterSymbols 函数中的相关说明。
		if (len(line) == 0) || (line[0] != '-') {
			continue
		}

		// 去掉 - 和空格
		line = bytes.TrimLeftFunc(line, func(r rune) bool {
			return r == '-' || unicode.IsSpace(r)
		})

		// 拿到首单词
		str := string(line)
		for index, c := range str {
			if unicode.IsSpace(c) {
				enum = append(enum, str[:index])
				continue LOOP
			}
		}
	}

	if len(enum) == 0 {
		return nil, nil
	}

	// 检查重复的枚举值
	sort.Strings(enum)
	for i := 1; i < len(enum); i++ {
		if enum[i-1] == enum[i] {
			return nil, errors.New("", "", 0, locale.ErrDuplicateValue)
		}
	}

	return convertEnumType(enum, typ)
}

func convertEnumType(enum []string, typ string) ([]interface{}, error) {
	fn := getConvertFunc(typ)

	ret := make([]interface{}, 0, len(enum))
	for _, elem := range enum {
		v, err := fn(elem)
		if err != nil {
			return nil, err
		}
		ret = append(ret, v)
	}

	return ret, nil
}

func getConvertFunc(typ string) convertFunc {
	fn, found := converters[typ]
	if !found {
		return stringConvert
	}
	return fn
}

// 类型转换函数定义
type convertFunc func(val string) (interface{}, error)

var (
	numberConvert = func(v string) (interface{}, error) {
		return strconv.ParseInt(v, 10, 64)
	}
	stringConvert = func(v string) (interface{}, error) {
		return v, nil
	}
	boolConvert = func(v string) (interface{}, error) {
		return strconv.ParseBool(v)
	}
	converters = map[string]convertFunc{
		Number:  numberConvert,
		Integer: numberConvert,
		String:  stringConvert,
		Bool:    boolConvert,
	}
)
