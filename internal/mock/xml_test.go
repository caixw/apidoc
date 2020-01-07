// SPDX-License-Identifier: MIT

package mock

import (
	"encoding/xml"
	"strings"
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v6/doc"
)

func TestValidXML(t *testing.T) {
	a := assert.New(t)

	for _, item := range data {
		err := validXML(item.Type, []byte(item.XML))
		a.NotError(err, "测试 %s 时返回错误 %s", item.Title, err)
	}

	p := &doc.Request{
		Name: "root",
		Type: doc.Object,
		Items: []*doc.Param{
			{
				Name: "id",
				Type: doc.Number,
				XML:  doc.XML{XMLAttr: true},
			},
			{
				Name: "desc",
				Type: doc.String,
				XML:  doc.XML{XMLExtract: true},
			},
		},
	}
	content := `<root id="1024"><desc>1024</desc></root>`
	a.Error(validXML(p, []byte(content)))
}

func TestBuildXML(t *testing.T) {
	a := assert.New(t)

	for _, item := range data {
		data, err := buildXML(item.Type)
		a.NotError(err, "测试 %s 返回了错误信息 %s", item.Title, err).
			Equal(string(data), item.XML, "测试 %s 返回的数据不相等 v1:%s,v2:%s", item.Title, string(data), item.XML)
	}
}

func TestXMLValidator_find(t *testing.T) {
	item := data[len(data)-1]

	a := assert.New(t)
	v := &xmlValidator{
		param:   item.Type.Param(),
		decoder: xml.NewDecoder(strings.NewReader(item.XML)),
	}

	v.names = []string{}
	p := v.find()
	a.Nil(p)

	v.names = nil
	p = v.find()
	a.Nil(p)

	v.names = []string{}
	p = v.find()
	a.Nil(p)

	v.names = []string{"root"}
	p = v.find()
	a.NotNil(p).Equal(p.Type, doc.Object)

	v.names = []string{"not-exists"}
	p = v.find()
	a.Nil(p)

	v.names = []string{"root", "group", "id"}
	p = v.find()
	a.NotNil(p).Equal(p.Type, doc.Number)

	v.names = []string{"root", "group", "tags", "id"}
	p = v.find()
	a.NotNil(p).Equal(p.Type, doc.Number)
}

func TestValidXMLParamValue(t *testing.T) {
	a := assert.New(t)

	// None
	a.NotError(validXMLParamValue(&doc.Param{}, "", ""))
	a.Error(validXMLParamValue(&doc.Param{}, "", "xx"))
	a.NotError(validXMLParamValue(&doc.Param{Type: doc.None}, "", ""))
	a.Error(validXMLParamValue(&doc.Param{Type: doc.None}, "", "xx"))

	// Number
	a.NotError(validXMLParamValue(&doc.Param{Type: doc.Number}, "", "1111"))
	a.NotError(validXMLParamValue(&doc.Param{Type: doc.Number}, "", "0"))
	a.NotError(validXMLParamValue(&doc.Param{Type: doc.Number}, "", "-11"))
	a.NotError(validXMLParamValue(&doc.Param{Type: doc.Number}, "", "-1024.11"))
	a.NotError(validXMLParamValue(&doc.Param{Type: doc.Number}, "", "-1024.1111234"))
	a.Error(validXMLParamValue(&doc.Param{Type: doc.Number}, "", "fxy0"))

	// String
	a.NotError(validXMLParamValue(&doc.Param{Type: doc.String}, "", "fxy0"))
	a.NotError(validXMLParamValue(&doc.Param{Type: doc.String}, "", ""))

	// Bool
	a.NotError(validXMLParamValue(&doc.Param{Type: doc.Bool}, "", "true"))
	a.NotError(validXMLParamValue(&doc.Param{Type: doc.Bool}, "", "false"))
	a.NotError(validXMLParamValue(&doc.Param{Type: doc.Bool}, "", "1"))
	a.Error(validXMLParamValue(&doc.Param{Type: doc.Bool}, "", "false/true"))

	// Other
	a.Error(validXMLParamValue(&doc.Param{Type: doc.Object}, "", ""))
	a.Error(validXMLParamValue(&doc.Param{Type: doc.Object}, "", "{}"))
	a.Error(validXMLParamValue(&doc.Param{Type: "xxx"}, "", "{}"))
	a.Error(validXMLParamValue(&doc.Param{Type: "xxx"}, "", ""))

	// bool enum
	p := &doc.Param{Type: doc.Bool, Enums: []*doc.Enum{
		{Value: "true"},
		{Value: "false"},
	}}
	a.NotError(validXMLParamValue(p, "", "true"))

	// 不存在于枚举
	p = &doc.Param{Type: doc.Bool, Enums: []*doc.Enum{
		{Value: "true"},
	}}
	a.Error(validXMLParamValue(p, "", "false"))

	// number enum
	p = &doc.Param{Type: doc.Number, Enums: []*doc.Enum{
		{Value: "1"},
		{Value: "-1.2"},
	}}
	a.NotError(validXMLParamValue(p, "", "1"))
	a.NotError(validXMLParamValue(p, "", "-1.2"))

	// 不存在于枚举
	p = &doc.Param{Type: doc.Number, Enums: []*doc.Enum{
		{Value: "1"},
		{Value: "-1.2"},
	}}
	a.Error(validXMLParamValue(p, "", "false"))
}

func TestGetXMLValue(t *testing.T) {
	a := assert.New(t)

	v, err := getXMLValue(&doc.Param{})
	a.NotError(err).Equal(v, "")

	v, err = getXMLValue(&doc.Param{Type: doc.None})
	a.NotError(err).Equal(v, "")

	v, err = getXMLValue(&doc.Param{Type: doc.Bool})
	a.NotError(err).Equal(v, true)

	v, err = getXMLValue(&doc.Param{Type: doc.Number})
	a.NotError(err).Equal(v, 1024)

	v, err = getXMLValue(&doc.Param{Type: doc.String})
	a.NotError(err).Equal(v, "1024")

	v, err = getXMLValue(&doc.Param{Type: doc.Object})
	a.Error(err).Nil(v)

	v, err = getXMLValue(&doc.Param{Type: "not-exists"})
	a.Error(err).Nil(v)
}
