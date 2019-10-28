// SPDX-License-Identifier: MIT

package doc

import (
	"encoding/xml"
	"fmt"
	"testing"

	"github.com/issue9/assert"
)

var (
	_ xml.Marshaler = Richtext{}
	_ fmt.Stringer  = Richtext{}
)

func TestRichtext_Marshal(t *testing.T) {
	a := assert.New(t)

	type Object struct {
		XMLName struct{} `xml:"xml"`
		Value   Richtext `xml:"value"`
	}

	obj := &Object{
		Value: Richtext{Text: "<a>test</a>"},
	}
	str := `<xml><value><a>test</a></value></xml>`

	data, err := xml.Marshal(obj)
	a.NotError(err).Equal(string(data), str)

	obj1 := &Object{}
	a.NotError(xml.Unmarshal([]byte(str), obj1))
	a.Equal(obj1, obj)

	// 空值
	obj = &Object{
		Value: Richtext{},
	}
	str = `<xml></xml>`

	data, err = xml.Marshal(obj)
	a.NotError(err).Equal(string(data), str)

	obj1 = &Object{}
	a.NotError(xml.Unmarshal([]byte(str), obj1))
	a.Equal(obj1, obj)
}

func TestShadowRichtext_Marshal(t *testing.T) {
	a := assert.New(t)

	type Object struct {
		XMLName struct{}       `xml:"xml"`
		Value   shadowRichtext `xml:"value"`
	}

	obj := &Object{
		Value: shadowRichtext{Text: "<p>cdata</p>"},
	}
	str := `<xml><value><p>cdata</p></value></xml>`

	data, err := xml.Marshal(obj)
	a.NotError(err).Equal(string(data), str)

	obj1 := &Object{}
	a.NotError(xml.Unmarshal([]byte(str), obj1))
	a.Equal(obj1, obj)

	// 空值
	obj = &Object{
		Value: shadowRichtext{},
	}
	str = `<xml><value></value></xml>`

	data, err = xml.Marshal(obj)
	a.NotError(err).Equal(string(data), str)

	obj1 = &Object{}
	a.NotError(xml.Unmarshal([]byte(str), obj1))
	a.Equal(obj1, obj)
}
