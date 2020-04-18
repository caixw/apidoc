// SPDX-License-Identifier: MIT

package token

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"reflect"
)

// Encoder 实现将内容转换成 XML 的接口
type Encoder interface {
	EncodeXML() (string, error)
}

// AttrEncoder 实现将内容作为 XML 属性的接口
type AttrEncoder interface {
	EncodeXMLAttr() (string, error)
}

var (
	attrEncoderType = reflect.TypeOf((*AttrEncoder)(nil)).Elem()
	encoderType     = reflect.TypeOf((*Encoder)(nil)).Elem()
)

// Encode 将 v 转换成 XML 内容
func Encode(indent, name string, v interface{}) ([]byte, error) {
	rv := reflect.ValueOf(v)
	if !rv.IsValid() {
		return nil, nil
	}

	buf := new(bytes.Buffer)
	e := xml.NewEncoder(buf)
	e.Indent("", indent)
	n := newNode(name, rv)

	if err := n.encode(e); err != nil {
		return nil, err
	}

	if err := e.Flush(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (n *node) encode(e *xml.Encoder) error {
	start, err := n.buildStartElement()
	if err != nil {
		return err
	}

	if n.cdata.IsValid() {
		chardata, err := getElementValue(n.cdata.Value)
		if err != nil {
			return err
		}
		return e.EncodeElement(struct {
			string `xml:",cdata"`
		}{chardata}, start)
	}

	if n.content.IsValid() {
		chardata, err := getElementValue(n.content.Value)
		if err != nil {
			return err
		}
		return e.EncodeElement(xml.CharData(chardata), start)
	}

	return n.encodeElems(e, start)
}

func (n *node) encodeElems(e *xml.Encoder, start xml.StartElement) (err error) {
	if err = e.EncodeToken(start); err != nil {
		return err
	}
	for _, v := range n.elems {
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
		if !found && isScalar(v.Value) {
			chardata = fmt.Sprint(v.Interface())
			found = true
		}

		if found {
			start := xml.StartElement{Name: xml.Name{Local: v.name}}
			if err := e.EncodeElement(xml.CharData(chardata), start); err != nil {
				return err
			}
			continue
		}

		if err := newNode(v.name, v.Value).encode(e); err != nil {
			return err
		}
	}
	return e.EncodeToken(xml.EndElement{Name: xml.Name{Local: n.name}})
}

func (n *node) buildStartElement() (xml.StartElement, error) {
	start := xml.StartElement{
		Name: xml.Name{Local: n.name},
		Attr: make([]xml.Attr, 0, len(n.attrs)),
	}

	for _, v := range n.attrs {
		val, err := getAttributeValue(v.Value)
		if err != nil {
			return xml.StartElement{}, err
		}

		start.Attr = append(start.Attr, xml.Attr{
			Name:  xml.Name{Local: v.name},
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

func getElementValue(elem reflect.Value) (string, error) {
	if elem.CanInterface() && elem.Type().Implements(encoderType) {
		return elem.Interface().(Encoder).EncodeXML()
	} else if elem.CanAddr() {
		pv := elem.Addr()
		if pv.CanInterface() && pv.Type().Implements(encoderType) {
			return pv.Interface().(Encoder).EncodeXML()
		}
	}

	return fmt.Sprint(elem.Interface()), nil
}
