// SPDX-License-Identifier: MIT

package token

import (
	"strconv"
	"testing"

	"github.com/issue9/assert"
)

type (
	attrEncodeObject struct {
		ID int
	}
	attrEncodeInt int

	encodeObject struct {
		ID int
	}
	encodeInt int
)

var (
	_ AttrEncoder = &attrEncodeObject{}
	_ AttrEncoder = attrEncodeInt(5)

	_ Encoder = &encodeObject{}
	_ Encoder = encodeInt(5)
)

func (o *attrEncodeObject) EncodeXMLAttr() (string, error) {
	return strconv.Itoa(o.ID + 1), nil
}

func (o attrEncodeInt) EncodeXMLAttr() (string, error) {
	return strconv.Itoa(int(o) + 1), nil
}

func (o *encodeObject) EncodeXML() (string, error) {
	return strconv.Itoa(o.ID + 1), nil
}

func (o encodeInt) EncodeXML() (string, error) {
	return strconv.Itoa(int(o) + 1), nil
}

func TestEncode(t *testing.T) {
	a := assert.New(t)

	type nestObject struct {
		ID   *encodeObject `apidoc:"id,elem"`
		Name string        `apidoc:"name,attr"`
	}

	data := []*struct {
		name   string
		object interface{}
		xml    string
		err    bool
	}{
		{},

		{
			name:   "apidoc",
			object: &struct{}{},
			xml:    "<apidoc></apidoc>",
		},

		{
			name: "apidoc",
			object: &struct {
				ID int `apidoc:"id,attr"`
			}{
				ID: 11,
			},
			xml: `<apidoc id="11"></apidoc>`,
		},

		{
			name: "apidoc",
			object: &struct {
				ID   int    `apidoc:"id,attr"`
				Name string `apidoc:",attr"`
			}{
				ID:   11,
				Name: "name",
			},
			xml: `<apidoc id="11" Name="name"></apidoc>`,
		},

		{
			name: "apidoc",
			object: &struct {
				ID   int    `apidoc:"id,attr"`
				Name string `apidoc:"name,elem"`
			}{
				ID:   11,
				Name: "name",
			},
			xml: `<apidoc id="11"><name>name</name></apidoc>`,
		},

		{
			name: "apidoc",
			object: &struct {
				ID    int    `apidoc:"id,attr"`
				CData string `apidoc:",cdata"`
			}{
				ID:    11,
				CData: "<h1>h1</h1>",
			},
			xml: `<apidoc id="11"><![CDATA[<h1>h1</h1>]]></apidoc>`,
		},

		{
			name: "apidoc",
			object: &struct {
				ID      int    `apidoc:"id,attr"`
				Content string `apidoc:",content"`
			}{
				ID:      11,
				Content: "<111",
			},
			xml: `<apidoc id="11">&lt;111</apidoc>`,
		},

		{
			name: "apidoc",
			object: &struct {
				ID           encodeInt         `apidoc:"id,attr"`
				IDObject     *encodeObject     `apidoc:"id_object,attr"`
				AttrID       attrEncodeInt     `apidoc:"attr_id,attr"`
				AttrIDObject *attrEncodeObject `apidoc:"attr_id_object,attr"`
			}{
				ID:           11,
				IDObject:     &encodeObject{ID: 11},
				AttrID:       11,
				AttrIDObject: &attrEncodeObject{ID: 11},
			},
			xml: `<apidoc id="11" id_object="{11}" attr_id="12" attr_id_object="12"></apidoc>`,
		},

		{
			name: "apidoc",
			object: &struct {
				ID           encodeInt         `apidoc:"id,elem"`
				IDObject     *encodeObject     `apidoc:"id_object,elem"`
				AttrID       attrEncodeInt     `apidoc:"attr_id,elem"`
				AttrIDObject *attrEncodeObject `apidoc:"attr_id_object,elem"`
			}{
				ID:           11,
				IDObject:     &encodeObject{ID: 11},
				AttrID:       11,
				AttrIDObject: &attrEncodeObject{ID: 11},
			},
			xml: `<apidoc><id>12</id><id_object>12</id_object><attr_id>11</attr_id><attr_id_object><ID>11</ID></attr_id_object></apidoc>`,
		},

		{ // 嵌套
			name: "apidoc",
			object: &struct {
				Object *nestObject `apidoc:"object,elem"`
			}{
				Object: &nestObject{
					ID:   &encodeObject{ID: 11},
					Name: "name",
				},
			},
			xml: `<apidoc><object name="name"><id>12</id></object></apidoc>`,
		},

		{ // cdata
			name: "apidoc",
			object: &struct {
				Cdata *encodeObject `apidoc:",cdata"`
			}{
				Cdata: &encodeObject{ID: 11},
			},
			xml: `<apidoc><![CDATA[12]]></apidoc>`,
		},

		{ // content
			name: "apidoc",
			object: &struct {
				Content encodeInt `apidoc:",content"`
			}{
				Content: 11,
			},
			xml: `<apidoc>12</apidoc>`,
		},
	}

	for i, item := range data {
		xml, err := Encode("", item.name, item.object)

		if item.err {
			a.Error(err, "not error at %d", i).
				Nil(xml, "not nil at %d", i)
			continue
		}

		a.NotError(err, "err %s at %d", err, i).
			Equal(string(xml), item.xml, "not equal at %d\nv1=%s\nv2=%s", i, string(xml), item.xml)
	}
}
