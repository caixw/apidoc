// SPDX-License-Identifier: MIT

package mock

import (
	"encoding/xml"
	"strings"
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v5/doc"
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
		param:   item.Type.ToParam(),
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
