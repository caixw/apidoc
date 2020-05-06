// SPDX-License-Identifier: MIT

package token

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"

	"github.com/caixw/apidoc/v7/core"
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
	RootName struct{}   `apidoc:"apidoc,attr,usage-root"`
	ID       intTest    `apidoc:"id,attr,usage"`
	Name     stringTest `apidoc:"name,elem,usage"`
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

	_ Sanitizer = &objectTest{}
)

func (i *intTest) DecodeXML(p *Parser, start *StartElement) (*EndElement, error) {
	for {
		t, err := p.Token()
		if err == io.EOF {
			return nil, nil
		} else if err != nil {
			return nil, err
		}

		switch elem := t.(type) {
		case *EndElement:
			if elem.Name.Value == start.Name.Value {
				return elem, nil
			}
		case *String:
			v, err := strconv.Atoi(strings.TrimSpace(elem.Value))
			if err != nil {
				return nil, err
			}
			i.Value = v
		default:
			panic(fmt.Sprintf("无效的类型 %s", reflect.TypeOf(t)))
		}
	}
}

func (i *intTest) DecodeXMLAttr(p *Parser, attr *Attribute) error {
	v, err := strconv.Atoi(strings.TrimSpace(attr.Value.Value))
	if err != nil {
		return p.WithError(attr.Value.Start, attr.Value.End, attr.Name.Value, err)
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

func (i *stringTest) DecodeXML(p *Parser, start *StartElement) (*EndElement, error) {
	for {
		t, err := p.Token()
		if err == io.EOF {
			return nil, nil
		} else if err != nil {
			return nil, err
		}

		switch elem := t.(type) {
		case *EndElement:
			if elem.Name.Value == start.Name.Value {
				return elem, nil
			}
		case *String:
			i.Value = elem.Value
		default:
			panic(fmt.Sprintf("无效的类型 %s", reflect.TypeOf(t)))
		}
	}
}

func (i *stringTest) DecodeXMLAttr(p *Parser, attr *Attribute) error {
	i.Value = attr.Value.Value
	return nil
}

func (i *stringTest) EncodeXML() (string, error) {
	return i.Value, nil
}

func (i *stringTest) EncodeXMLAttr() (string, error) {
	return i.Value, nil
}

func (o *objectTest) Sanitize(p *Parser) error {
	o.ID.Value++
	return nil
}

func (t *errIntTest) DecodeXML(p *Parser, start *StartElement) (core.Position, error) {
	return core.Position{}, errors.New("Decoder.DecodeXML")
}

func (t *errIntTest) DecodeXMLAttr(p *Parser, attr *Attribute) error {
	return errors.New("AttrDecoder.DecodeXMLAttr")
}

func (t *errIntTest) EncodeXML() (string, error) {
	return "", errors.New("Encoder.EncodeXML")
}

func (t *errIntTest) EncodeXMLAttr() (string, error) {
	return "", errors.New("AttrEncoder.EncodeXMLAttr")
}
