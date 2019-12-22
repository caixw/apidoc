// SPDX-License-Identifier: MIT

package mock

import (
	"bytes"
	"encoding/xml"
	"io"
	"strconv"
	"strings"

	"github.com/issue9/is"

	"github.com/caixw/apidoc/v5/doc"
	"github.com/caixw/apidoc/v5/internal/locale"
	"github.com/caixw/apidoc/v5/message"
)

type xmlValidator struct {
	param   *doc.Param
	decoder *xml.Decoder
	names   []string // 按顺序保存变量名称
}

func validXML(p *doc.Request, content []byte) error {
	if len(content) == 0 {
		if p == nil || p.Type == doc.None {
			return nil
		}
		return message.NewLocaleError("", "", 0, locale.ErrInvalidFormat)
	}

	validator := &xmlValidator{
		param:   p.ToParam(),
		decoder: xml.NewDecoder(bytes.NewReader(content)),
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

		switch v := token.(type) {
		case xml.StartElement:
			validator.pushName(v.Name.Local)
			for _, attr := range v.Attr {
				if err := validator.validValue(attr.Name.Local, attr.Value); err != nil {
					return err
				}
			}
		case xml.EndElement:
			validator.popName()
		case xml.CharData:
			if err := validator.validValue("", string(v)); err != nil {
				return err
			}
		case xml.Comment, xml.ProcInst, xml.Directive:
		}
	}
}

// 如果 t == "" 表示不需要验证类型，比如 null 可以赋值给任何类型
func (validator *xmlValidator) validValue(k, v string) error {
	field := strings.Join(validator.names, "/")

	p := validator.find(k)
	if p == nil {
		if k != "" {
			field += "/" + k
		}
		return message.NewLocaleError("", field, 0, locale.ErrNotFound)
	}

	if p.XMLAttr {
		field += "/@" + k
	}

	switch p.Type {
	case doc.Number:
		if !is.Number(v) {
			return message.NewLocaleError("", field, 0, locale.ErrInvalidFormat)
		}
	case doc.Bool:
		if _, err := strconv.ParseBool(v); err != nil {
			return message.NewLocaleError("", field, 0, locale.ErrInvalidFormat)
		}
	case doc.String:
		return nil
	case doc.None:
		if v != "" {
			return message.NewLocaleError("", field, 0, locale.ErrInvalidValue)
		}
	case doc.Object:
	}

	if p.IsEnum() {
		for _, enum := range p.Enums {
			if enum.Value == v {
				return nil
			}
		}
		return message.NewLocaleError("", field, 0, locale.ErrInvalidValue)
	}

	return nil
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
//
// name 会被附加在 validator.names 之后，进行查询，如果为空，则只查询 validator.names 的值。
func (validator *xmlValidator) find(name string) *doc.Param {
	names := make([]string, len(validator.names))
	copy(names, validator.names)
	if name != "" {
		names = append(names, name)
	}

	if len(names) == 0 || validator.param == nil || names[0] != validator.param.Name {
		return nil
	}

	p := validator.param
	for _, name := range names[1:] {
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
