// SPDX-License-Identifier: MIT

package token

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"

	"github.com/caixw/apidoc/v6/core"
)

type anonymous struct {
	Attr1 intTest `apidoc:"attr1,attr,usage"`
	Elem1 intTest `apidoc:"elem1,elem,usage"`
}

type intTest struct {
	Base
	Value int
}

type stringTest struct {
	Base
	Value string
}

type errIntTest struct {
	Base
	Value int
}

// NOTE: objectTest 作为普通对象嵌套了 Decoder 等实例，本身不能实现这些接口。
type objectTest struct {
	Base
	ID   intTest    `apidoc:"id,attr,usage"`
	Name stringTest `apidoc:"name,elem,usage"`
}

var (
	_ AttrEncoder = &intTest{}
	_ AttrDecoder = &intTest{}
	_ Encoder     = &intTest{}
	_ Decoder     = &intTest{}

	_ AttrEncoder = &stringTest{}
	_ AttrDecoder = &stringTest{}
	_ Encoder     = &stringTest{}
	_ Decoder     = &stringTest{}
)

func (i *intTest) DecodeXML(p *Parser, start *StartElement) (core.Position, error) {
	for {
		t, err := p.Token()
		if err != nil {
			return core.Position{}, err
		}
		if t == nil {
			return core.Position{}, io.EOF
		}

		switch elem := t.(type) {
		case *EndElement:
			if elem.Name.Value == start.Name.Value {
				return elem.End, nil
			}
		case *String:
			v, err := strconv.Atoi(strings.TrimSpace(elem.Value))
			if err != nil {
				return core.Position{}, err
			}
			i.Value = v
		default:
			panic(fmt.Sprintf("无效的类型 %s", reflect.TypeOf(t)))
		}
	}
}

func (i *intTest) DecodeXMLAttr(attr *Attribute) error {
	v, err := strconv.Atoi(strings.TrimSpace(attr.Value.Value))
	if err != nil {
		return err
	}
	i.Value = v
	return nil
}

func (i *intTest) EncodeXML() (string, error) {
	return strconv.Itoa(i.Value), nil
}

func (i *intTest) EncodeXMLAttr() (string, error) {
	return strconv.Itoa(i.Value), nil
}

func (i *stringTest) DecodeXML(p *Parser, start *StartElement) (core.Position, error) {
	for {
		t, err := p.Token()
		if err != nil {
			return core.Position{}, err
		}
		if t == nil {
			return core.Position{}, io.EOF
		}

		switch elem := t.(type) {
		case *EndElement:
			if elem.Name.Value == start.Name.Value {
				return elem.End, nil
			}
		case *String:
			i.Value = elem.Value
		default:
			panic(fmt.Sprintf("无效的类型 %s", reflect.TypeOf(t)))
		}
	}
}

func (i *stringTest) DecodeXMLAttr(attr *Attribute) error {
	i.Value = attr.Value.Value
	return nil
}

func (i *stringTest) EncodeXML() (string, error) {
	return i.Value, nil
}

func (i *stringTest) EncodeXMLAttr() (string, error) {
	return i.Value, nil
}

func (t *errIntTest) DecodeXML(p *Parser, start *StartElement) (core.Position, error) {
	return core.Position{}, errors.New("Decoder.DecodeXML")
}

func (t *errIntTest) DecodeXMLAttr(attr *Attribute) error {
	return errors.New("AttrDecoder.DecodeXMLAttr")
}

func (t *errIntTest) EncodeXML() (string, error) {
	return "", errors.New("Encoder.EncodeXML")
}

func (t *errIntTest) EncodeXMLAttr() (string, error) {
	return "", errors.New("AttrEncoder.EncodeXMLAttr")
}
