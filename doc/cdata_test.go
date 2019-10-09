// SPDX-License-Identifier: MIT

package doc

import (
	"encoding/xml"
	"fmt"
	"testing"

	"github.com/issue9/assert"
)

var (
	_ xml.Marshaler = CDATA{}
	_ fmt.Stringer  = CDATA{}
)

func TestCDATA_Marshal(t *testing.T) {
	a := assert.New(t)

	type Object struct {
		XMLName struct{} `xml:"xml"`
		Value   CDATA    `xml:"value"`
	}

	obj := &Object{
		Value: CDATA{Text: "cdata"},
	}
	str := `<xml><value><![CDATA[cdata]]></value></xml>`

	data, err := xml.Marshal(obj)
	a.NotError(err).Equal(string(data), str)

	obj1 := &Object{}
	a.NotError(xml.Unmarshal([]byte(str), obj1))
	a.Equal(obj1, obj)

	// 空值
	obj = &Object{
		Value: CDATA{},
	}
	str = `<xml></xml>`

	data, err = xml.Marshal(obj)
	a.NotError(err).Equal(string(data), str)

	obj1 = &Object{}
	a.NotError(xml.Unmarshal([]byte(str), obj1))
	a.Equal(obj1, obj)
}

func TestShadowCDATA_Marshal(t *testing.T) {
	a := assert.New(t)

	type Object struct {
		XMLName struct{}    `xml:"xml"`
		Value   shadowCDATA `xml:"value"`
	}

	obj := &Object{
		Value: shadowCDATA{Text: "cdata"},
	}
	str := `<xml><value><![CDATA[cdata]]></value></xml>`

	data, err := xml.Marshal(obj)
	a.NotError(err).Equal(string(data), str)

	obj1 := &Object{}
	a.NotError(xml.Unmarshal([]byte(str), obj1))
	a.Equal(obj1, obj)

	// 空值
	obj = &Object{
		Value: shadowCDATA{},
	}
	str = `<xml><value></value></xml>`

	data, err = xml.Marshal(obj)
	a.NotError(err).Equal(string(data), str)

	obj1 = &Object{}
	a.NotError(xml.Unmarshal([]byte(str), obj1))
	a.Equal(obj1, obj)
}
