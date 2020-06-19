// SPDX-License-Identifier: MIT

package mock

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/issue9/errwrap"
	"github.com/issue9/is"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/ast"
	"github.com/caixw/apidoc/v7/internal/locale"
)

type jsonValidator struct {
	param *ast.Param

	// 按顺序表示的状态
	// 可以是 [ 表示在数组中，{ 表示在对象，: 表示下一个值必须是属性，空格表示其它状态
	states []byte

	names []string // 按顺序保存变量名称
}

func validJSON(p *ast.Request, content []byte) error {
	if p == nil {
		if bytes.Equal(content, []byte("null")) {
			return nil
		}
		return core.NewSyntaxError(core.Location{}, "", locale.ErrInvalidFormat)
	} else if p.Type.V() == ast.TypeNone && len(content) == 0 {
		return nil
	}

	if !json.Valid(content) {
		return core.NewSyntaxError(core.Location{}, "", locale.ErrInvalidFormat)
	}

	validator := &jsonValidator{
		param:  p.Param(),
		states: []byte{0}, // 状态有默认值
		names:  []string{},
	}

	return validator.valid(json.NewDecoder(bytes.NewReader(content)))
}

func (validator *jsonValidator) valid(d *json.Decoder) error {
	for {
		token, err := d.Token()
		if errors.Is(err, io.EOF) && token == nil { // 正常结束
			return nil
		}
		if err != nil {
			return err
		}
		if token == nil { // 对应 JSON null
			return nil
		}

		switch v := token.(type) {
		case string: // json string
			switch validator.state() {
			case ':': // 字符串类型的值
				err = validator.validValue(ast.TypeString, v)
				validator.popState()
				validator.popName()
			case '[':
				err = validator.validValue(ast.TypeString, v)
			default: // 属性名
				validator.pushState(':')
				validator.pushName(v)
			}

			if err != nil {
				return err
			}
		case json.Delim: // [、]、{、}
			switch v {
			case '[':
				validator.pushState('[')
			case ']':
				validator.popName()

				validator.popState()
				if validator.state() == ':' { // {xx: [] } 类似这种格式，需要同时弹出两个状态
					validator.popState()
				}
			case '{':
				validator.pushState('{')
			case '}':
				validator.popName()

				validator.popState()
				if validator.state() == ':' {
					validator.popState()
				}
			}
		case bool: // json bool
			err = validator.validValue(ast.TypeBool, v)
			if validator.state() != '[' {
				validator.popState()
				validator.popName()
			}
		case float64, json.Number: // json number
			err = validator.validValue(ast.TypeNumber, v)
			if validator.state() != '[' { // 只有键值对结束时，才弹出键名
				validator.popState()
				validator.popName()
			}
		}

		if err != nil {
			return err
		}
	}
}

var jsonNumberType = reflect.TypeOf(json.Number("1"))

// 如果 t == "" 表示不需要验证类型，比如 null 可以赋值给任何类型
func (validator *jsonValidator) validValue(t string, v interface{}) error {
	field := strings.Join(validator.names, ".")

	p := validator.find()
	if p == nil {
		return core.NewSyntaxError(core.Location{}, field, locale.ErrNotFound)
	}

	pt := p.Type.V()

	if primitive, _ := ast.ParseType(pt); primitive != t {
		return core.NewSyntaxError(core.Location{}, field, locale.ErrInvalidFormat)
	}

	switch pt {
	case ast.TypeEmail:
		if !is.Email(v) {
			return core.NewSyntaxError(core.Location{}, field, locale.ErrInvalidFormat)
		}
	case ast.TypeURL:
		if !is.URL(v) {
			return core.NewSyntaxError(core.Location{}, field, locale.ErrInvalidFormat)
		}
	case ast.TypeInt, ast.TypeFloat: // 数值类型都被 json 解释为 float64，无法判断值是浮点还是整数。
	}

	if isEnum(p) {
		for _, enum := range p.Enums {
			if enum.Value.V() == fmt.Sprint(v) {
				return nil
			}
		}
		return core.NewSyntaxError(core.Location{}, field, locale.ErrInvalidValue)
	}

	return nil
}

// 返回当前的状态
func (validator *jsonValidator) state() byte {
	if len(validator.states) > 0 {
		return validator.states[len(validator.states)-1]
	}
	return 0
}

func (validator *jsonValidator) pushState(state byte) {
	validator.states = append(validator.states, state)
}

func (validator *jsonValidator) popState() {
	if len(validator.states) > 0 {
		validator.states = validator.states[:len(validator.states)-1]
	}
}

func (validator *jsonValidator) pushName(name string) {
	validator.names = append(validator.names, name)
}

func (validator *jsonValidator) popName() {
	if len(validator.names) > 0 {
		validator.names = validator.names[:len(validator.names)-1]
	}
}

// 如果 names 为空，返回 validator.param
func (validator *jsonValidator) find() *ast.Param {
	p := validator.param

LOOP:
	for _, name := range validator.names {
		for _, pp := range p.Items {
			if pp.Name.V() == name {
				p = pp
				continue LOOP
			}
		}
		return nil
	}

	return p
}

type jsonBuilder struct {
	w      *errwrap.Buffer
	deep   int
	indent string // 单次的缩进
}

func buildJSON(p *ast.Request, indent string, g *GenOptions) ([]byte, error) {
	if p != nil && p.Type.V() == ast.TypeNone {
		return nil, nil
	}

	builder := &jsonBuilder{
		w:      &errwrap.Buffer{},
		indent: indent,
	}

	if err := builder.encode(p.Param(), true, g); err != nil {
		return nil, err
	}

	return builder.w.Bytes(), nil
}

func (builder *jsonBuilder) encode(p *ast.Param, chkArray bool, g *GenOptions) error {
	if p == nil {
		return builder.writeValue(nil).w.Err
	}

	if p.Array.V() && chkArray {
		builder.w.WString("[\n")
		builder.deep++

		size := g.generateSliceSize()
		last := size - 1
		for i := 0; i < size; i++ {
			if err := builder.writeIndent().encode(p, false, g); err != nil {
				return err
			}

			if i < last {
				builder.w.WString(",\n")
			} else {
				builder.w.WString("\n")
			}
		}

		builder.deep--
		return builder.writeIndent().w.WString("]").Err
	}

	switch primitive, _ := ast.ParseType(p.Type.V()); primitive {
	case ast.TypeNone:
		builder.writeValue(nil)
	case ast.TypeBool:
		builder.writeValue(g.generateBool())
	case ast.TypeNumber:
		builder.writeValue(g.generateNumber(p))
	case ast.TypeString:
		builder.writeValue(g.generateString(p))
	case ast.TypeObject:
		builder.w.WString("{\n")
		builder.deep++

		last := len(p.Items) - 1
		for index, item := range p.Items {
			builder.writeIndent().w.WString(`"`).WString(item.Name.V()).WString(`"`).WString(": ")

			if err := builder.encode(item, true, g); err != nil {
				return err
			}

			if index < last {
				builder.w.WString(",\n")
			} else {
				builder.w.WString("\n")
			}
		}

		builder.deep--
		builder.writeIndent().w.WString("}")
	}

	return builder.w.Err
}

func (builder *jsonBuilder) writeIndent() *jsonBuilder {
	builder.w.WString(strings.Repeat(builder.indent, builder.deep))
	return builder
}

// v 只能是基本类型
func (builder *jsonBuilder) writeValue(v interface{}) *jsonBuilder {
	if builder.w.Err != nil {
		return builder
	}

	vv, err := json.Marshal(v)
	if err != nil {
		builder.w.Err = err
		return builder
	}

	builder.w.WBytes(vv)
	return builder
}
