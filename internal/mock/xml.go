// SPDX-License-Identifier: MIT

package mock

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"strings"

	"github.com/caixw/apidoc/v5/doc"
	"github.com/caixw/apidoc/v5/internal/locale"
	"github.com/caixw/apidoc/v5/message"
)

type xmlValidator struct {
	param   *doc.Param
	decoder *xml.Decoder

	// 按顺序表示的状态
	// 可以是 [ 表示在数组中，{ 表示在对象，: 表示下一个值必须是属性，空格表示其它状态
	states []byte

	// 按顺序保存变量名称
	names []string
}

func validXML(p *doc.Request, content []byte) error {
	if p == nil && bytes.Equal(content, []byte("null")) {
		return nil
	}

	if (p.Type == doc.None || p.Type == "") && len(content) == 0 {
		return nil
	}

	validator := &xmlValidator{
		param:   p.ToParam(),
		decoder: xml.NewDecoder(bytes.NewReader(content)),
		states:  []byte{' '}, // 状态有默认值
		names:   []string{},
	}

	return validator.valid()
}

func (validator *xmlValidator) valid() error {
	for {
		token, err := validator.decoder.Token()
		if err == io.EOF && token == nil { // 正常结束
			return nil
		}

		if err != nil {
			return err
		}

		switch token.(type) {
		case xml.StartElement:
			// TODO
		case xml.EndElement:
		case xml.CharData:
		case xml.Comment:
		case xml.ProcInst:
		case xml.Directive:
		}

		if err != nil {
			return err
		}
	}
}

// 如果 t == "" 表示不需要验证类型，比如 null 可以赋值给任何类型
func (validator *xmlValidator) validValue(t doc.Type, v interface{}) error {
	field := strings.Join(validator.names, ".")

	p := validator.find()
	if p == nil {
		return message.NewLocaleError("", field, 0, locale.ErrNotFound)
	}

	if t == "" {
		return nil
	}

	if p.Type != t {
		return message.NewLocaleError("", field, 0, locale.ErrInvalidFormat)
	}

	if p.IsEnum() {
		for _, enum := range p.Enums {
			if enum.Value == fmt.Sprint(v) {
				return nil
			}
		}
		return message.NewLocaleError("", field, 0, locale.ErrInvalidValue)
	}

	return nil
}

// 返回当前的状态
func (validator *xmlValidator) state() byte {
	return validator.states[len(validator.states)-1]
}

func (validator *xmlValidator) pushState(state byte) *xmlValidator {
	validator.states = append(validator.states, state)
	return validator
}

func (validator *xmlValidator) popState() *xmlValidator {
	if len(validator.states) > 0 {
		validator.states = validator.states[:len(validator.states)-1]
	}
	return validator
}

func (validator *xmlValidator) pushName(name string) *xmlValidator {
	validator.names = append(validator.names, name)
	return validator
}

func (validator *xmlValidator) popName() *xmlValidator {
	if len(validator.names) > 0 {
		validator.names = validator.names[:len(validator.names)-1]
	}
	return validator
}

// 如果 names 为空，返回 validator.param
func (validator *xmlValidator) find() *doc.Param {
	p := validator.param
	for _, name := range validator.names {
		found := false
		for _, pp := range p.Items {
			if pp.Name == name {
				p = pp
				found = true
				break
			}
		}

		if !found {
			return nil
		}
	}

	return p
}

func buildXML(p *doc.Request) ([]byte, error) {
	// TODO
	return nil, nil
}
