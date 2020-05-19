// SPDX-License-Identifier: MIT

package token

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"reflect"

	"github.com/issue9/is"

	"github.com/caixw/apidoc/v7/internal/node"
)

// Encoder 将元素内容编码成 XML 内容
type Encoder interface {
	// 仅需要返回元素内容的 XML 编码，不需要包含本身的标签和属性。
	EncodeXML() (string, error)
}

// AttrEncoder 将属性值编码成符合 XML 规范的值
type AttrEncoder interface {
	// 仅需要返回属性的 XML 表示，不需要包含属性值的引号字符。
	EncodeXMLAttr() (string, error)
}

var (
	attrEncoderType = reflect.TypeOf((*AttrEncoder)(nil)).Elem()
	encoderType     = reflect.TypeOf((*Encoder)(nil)).Elem()
)

// Encode 将 v 转换成 XML 内容
func Encode(indent string, v interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	e := xml.NewEncoder(buf)
	e.Indent("", indent)

	if err := encode(node.New("", reflect.ValueOf(v)), e); err != nil {
		return nil, err
	}

	if err := e.Flush(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func encode(n *node.Node, e *xml.Encoder) error {
	start, err := buildStartElement(n)
	if err != nil {
		return err
	}

	if n.CData != nil && !isOmitempty(n.CData) {
		chardata, err := getContentValue(n.CData.Value)
		if err != nil {
			return err
		}

		return e.EncodeElement(struct {
			string `xml:",cdata"`
		}{chardata}, start)
	}

	if n.Content != nil && !isOmitempty(n.Content) {
		chardata, err := getContentValue(n.Content.Value)
		if err != nil {
			return err
		}
		return e.EncodeElement(xml.CharData(chardata), start)
	}

	return encodeElements(n, e, start)
}

func encodeElements(n *node.Node, e *xml.Encoder, start xml.StartElement) (err error) {
	if err = e.EncodeToken(start); err != nil {
		return err
	}
	for _, v := range n.Elements {
		if err := encodeElement(e, v); err != nil {
			return err
		}
	}

	return e.EncodeToken(xml.EndElement{Name: xml.Name{Local: n.Value.Name}})
}

func encodeElement(e *xml.Encoder, v *node.Value) (err error) {
	if isOmitempty(v) {
		return nil
	}

	var chardata string
	var found bool

	if v.CanInterface() && v.Type().Implements(encoderType) {
		if chardata, err = v.Interface().(Encoder).EncodeXML(); err != nil {
			return err
		}
		found = true
	} else if !found && v.CanAddr() {
		pv := v.Addr()
		if pv.CanInterface() && pv.Type().Implements(encoderType) {
			if chardata, err = pv.Interface().(Encoder).EncodeXML(); err != nil {
				return err
			}
			found = true
		}
	}
	if !found && node.IsPrimitive(v.Value) {
		chardata = fmt.Sprint(v.Interface())
		found = true
	}

	if found {
		start := xml.StartElement{Name: xml.Name{Local: v.Name}}
		if err := e.EncodeElement(xml.CharData(chardata), start); err != nil {
			return err
		}
		return nil
	}

	if v.Kind() == reflect.Array || v.Kind() == reflect.Slice {
		for i := 0; i < v.Len(); i++ {
			if err := encodeElement(e, node.NewValue(v.Name, v.Index(i), v.Omitempty, v.Usage)); err != nil {
				return err
			}
		}
		return nil
	}

	return encode(node.New(v.Name, v.Value), e)
}

func buildStartElement(n *node.Node) (xml.StartElement, error) {
	start := xml.StartElement{
		Name: xml.Name{Local: n.Value.Name},
		Attr: make([]xml.Attr, 0, len(n.Attributes)),
	}

	for _, v := range n.Attributes {
		if isOmitempty(v) {
			continue
		}

		val, err := getAttributeValue(v.Value)
		if err != nil {
			return xml.StartElement{}, err
		}

		start.Attr = append(start.Attr, xml.Attr{
			Name:  xml.Name{Local: v.Name},
			Value: val,
		})
	}

	return start, nil
}

func getAttributeValue(elem reflect.Value) (string, error) {
	if elem.CanInterface() && elem.Type().Implements(attrEncoderType) {
		return elem.Interface().(AttrEncoder).EncodeXMLAttr()
	} else if elem.CanAddr() {
		pv := elem.Addr()
		if pv.CanInterface() && pv.Type().Implements(attrEncoderType) {
			return pv.Interface().(AttrEncoder).EncodeXMLAttr()
		}
	}

	return fmt.Sprint(elem.Interface()), nil
}

// 获取 cdata 和 content 节点的的内容
func getContentValue(elem reflect.Value) (string, error) {
	elem = node.GetRealValue(elem)
	if elem.CanInterface() && elem.Type().Implements(encoderType) {
		return elem.Interface().(Encoder).EncodeXML()
	} else if elem.CanAddr() {
		if pv := elem.Addr(); pv.CanInterface() && pv.Type().Implements(encoderType) {
			return pv.Interface().(Encoder).EncodeXML()
		}
	}

	return fmt.Sprint(elem.Interface()), nil
}

func isOmitempty(v *node.Value) bool {
	return v.Omitempty && is.Empty(v.Value.Interface(), true)
}
