// SPDX-License-Identifier: MIT

package token

import (
	"reflect"
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/node"
)

var (
	_ Encoder = &CData{}
	_ Encoder = &String{}
)

func (cdata *CData) EncodeXML() (string, error) {
	return cdata.Value.Value, nil
}

func (s *String) EncodeXML() (string, error) {
	return s.Value, nil
}

func TestEncode(t *testing.T) {
	a := assert.New(t)

	type nestObject struct {
		ID   *intTag    `apidoc:"id,elem,usage,omitempty"`
		Name stringAttr `apidoc:"name,attr,usage,omitempty"`
	}

	data := []*struct {
		object            interface{}
		namespace, prefix string
		xml               string
		err               bool
	}{
		{
			object: &struct {
				RootName string `apidoc:"apidoc,meta,usage-apidoc"`
			}{},
			xml: "<apidoc></apidoc>",
		},

		{
			object: &struct {
				RootName string `apidoc:"apidoc,meta,usage-apidoc"`
			}{},
			namespace: core.XMLNamespace,
			xml:       `<apidoc xmlns="` + core.XMLNamespace + `"></apidoc>`,
		},

		{
			object: &struct {
				RootName string `apidoc:"apidoc,meta,usage-apidoc"`
			}{},
			prefix:    "aa",
			namespace: core.XMLNamespace,
			xml:       `<aa:apidoc xmlns:aa="` + core.XMLNamespace + `"></aa:apidoc>`,
		},

		{
			object: &struct {
				RootName string  `apidoc:"apidoc,meta,usage-apidoc"`
				ID       intAttr `apidoc:"id,attr,usage"`
			}{
				ID: intAttr{Value: 11},
			},
			namespace: core.XMLNamespace,
			prefix:    "bb",
			xml:       `<bb:apidoc bb:id="11" xmlns:bb="` + core.XMLNamespace + `"></bb:apidoc>`,
		},

		{
			object: &struct {
				RootName string  `apidoc:"apidoc,meta,usage-apidoc"`
				ID       intAttr `apidoc:"id,attr,usage"`
			}{
				ID: intAttr{Value: 11},
			},
			xml: `<apidoc id="11"></apidoc>`,
		},

		{ // 非 omitempty 属性，必须带上零值
			object: &struct {
				RootName string  `apidoc:"apidoc,meta,usage-apidoc"`
				ID       intAttr `apidoc:"id,attr,usage"`
			}{},
			xml: `<apidoc id="0"></apidoc>`,
		},

		{ // omitempty
			object: &struct {
				RootName string  `apidoc:"apidoc,meta,usage-apidoc"`
				ID       intAttr `apidoc:"id,attr,usage,omitempty"`
			}{},
			xml: `<apidoc></apidoc>`,
		},

		{ // omitempty
			object: &struct {
				RootName string   `apidoc:"apidoc,meta,usage-apidoc"`
				ID       *intAttr `apidoc:"id,attr,usage,omitempty"`
			}{},
			xml: `<apidoc></apidoc>`,
		},

		{
			object: &struct {
				RootName string      `apidoc:"root,meta,usage-apidoc"`
				ID       intAttr     `apidoc:"id,attr,usage"`
				Name     *stringAttr `apidoc:",attr,usage"`
			}{
				ID:   intAttr{Value: 11},
				Name: &stringAttr{Value: "name"},
			},
			xml: `<root id="11" Name="name"></root>`,
		},

		{ // 数组
			object: &struct {
				RootName string       `apidoc:"apidoc,meta,usage-apidoc"`
				ID       []intTag     `apidoc:"id,elem,usage"`
				Name     []*stringTag `apidoc:",elem,usage"`
			}{
				ID:   []intTag{{Value: 11}, {Value: 12}},
				Name: []*stringTag{{Value: "name1"}, {Value: "name2"}},
			},
			xml: `<apidoc><id>11</id><id>12</id><Name>name1</Name><Name>name2</Name></apidoc>`,
		},

		{
			object: &struct {
				RootName string    `apidoc:"apidoc,meta,usage-apidoc"`
				ID       *intAttr  `apidoc:"id,attr,usage"`
				Name     stringTag `apidoc:"name,elem,usage"`
			}{
				ID:   &intAttr{Value: 11},
				Name: stringTag{Value: "name"},
			},
			xml: `<apidoc id="11"><name>name</name></apidoc>`,
		},

		{
			object: &struct {
				RootName string  `apidoc:"apidoc,meta,usage-apidoc"`
				ID       intAttr `apidoc:"id,attr,usage"`
				CData    CData   `apidoc:",cdata,"`
			}{
				ID:    intAttr{Value: 11},
				CData: CData{Value: String{Value: "<h1>h1</h1>"}},
			},
			xml: `<apidoc id="11"><![CDATA[<h1>h1</h1>]]></apidoc>`,
		},

		{
			object: &struct {
				RootName string  `apidoc:"apidoc,meta,usage-apidoc"`
				ID       int     `apidoc:"id,attr,usage"`
				Content  *String `apidoc:",content"`
			}{
				ID:      11,
				Content: &String{Value: "<111"},
			},
			xml: `<apidoc id="11">&lt;111</apidoc>`,
		},

		{
			object: &struct {
				RootName string  `apidoc:"apidoc,meta,usage-apidoc"`
				ID       int     `apidoc:"id,attr,usage"`
				Content  *String `apidoc:",content"`
			}{
				ID:      11,
				Content: &String{Value: "<111"},
			},
			namespace: "urn",
			prefix:    "p",
			xml:       `<p:apidoc p:id="11" xmlns:p="urn">&lt;111</p:apidoc>`,
		},

		{
			object: &struct {
				RootName string `apidoc:"apidoc,meta,usage-apidoc"`
				ID       int    `apidoc:"id,attr,usage"`
				Content  string `apidoc:",content"`
			}{
				ID:      11,
				Content: "<111",
			},
			xml: `<apidoc id="11">&lt;111</apidoc>`,
		},

		{ // 嵌套
			object: &struct {
				RootName string      `apidoc:"apidoc,meta,usage-apidoc"`
				Object   *nestObject `apidoc:"object,elem,usage"`
			}{
				Object: &nestObject{
					ID:   &intTag{Value: 12},
					Name: stringAttr{Value: "name"},
				},
			},
			xml: `<apidoc><object name="name"><id>12</id></object></apidoc>`,
		},

		{ // 嵌套 cdata
			object: &struct {
				RootName string `apidoc:"apidoc,meta,usage-apidoc"`
				Cdata    *CData `apidoc:",cdata"`
			}{
				Cdata: &CData{Value: String{Value: "12"}},
			},
			xml: `<apidoc><![CDATA[12]]></apidoc>`,
		},

		{ // 嵌套 content
			object: &struct {
				RootName string  `apidoc:"apidoc,meta,usage-apidoc"`
				Content  *String `apidoc:",content"`
			}{
				Content: &String{Value: "11"},
			},
			xml: `<apidoc>11</apidoc>`,
		},

		{ // 嵌套，omitempty 属性
			object: &struct {
				RootName string      `apidoc:"apidoc,meta,usage-apidoc"`
				Object   *nestObject `apidoc:"object,elem,usage,omitempty"`
			}{},
			xml: `<apidoc></apidoc>`,
		},

		{ // 嵌套，omitempty 属性
			object: &struct {
				RootName string      `apidoc:"apidoc,meta,usage-apidoc"`
				Object   *nestObject `apidoc:"object,elem,usage,omitempty"`
			}{
				Object: &nestObject{
					ID: &intTag{Value: 12},
				},
			},
			xml: `<apidoc><object><id>12</id></object></apidoc>`,
		},

		{ // 嵌套，omitempty 属性，namespace
			object: &struct {
				RootName string      `apidoc:"apidoc,meta,usage-apidoc"`
				Object   *nestObject `apidoc:"object,elem,usage,omitempty"`
			}{
				Object: &nestObject{
					ID: &intTag{Value: 12},
				},
			},
			namespace: "urn",
			prefix:    "p",
			xml:       `<p:apidoc xmlns:p="urn"><p:object><p:id>12</p:id></p:object></p:apidoc>`,
		},

		{ // 嵌套，数组，omitempty 属性
			object: &struct {
				RootName string        `apidoc:"aa,meta,usage-apidoc"`
				Object   []*nestObject `apidoc:"object,elem,usage,omitempty"`
			}{
				Object: []*nestObject{
					{ID: &intTag{Value: 12}},
					{ID: &intTag{Value: 22}},
				},
			},
			xml: `<aa><object><id>12</id></object><object><id>22</id></object></aa>`,
		},

		{ // 嵌套，omitempty 属性
			object: &struct {
				RootName string      `apidoc:"apidoc,meta,usage-apidoc"`
				Object   *nestObject `apidoc:"object,elem,usage"`
			}{
				Object: &nestObject{},
			},
			xml: `<apidoc><object></object></apidoc>`,
		},
	}

	for i, item := range data {
		xml, err := Encode("", item.object, item.namespace, item.prefix)

		if item.err {
			a.Error(err, "not error at %d", i).
				Nil(xml, "not nil at %d", i)
			continue
		}

		a.NotError(err, "err %s at %d", err, i).
			Equal(string(xml), item.xml, "not equal at %d\nv1=%s\nv2=%s", i, string(xml), item.xml)
	}
}

func TestNode_isOmitempty(t *testing.T) {
	a := assert.New(t)

	v := &node.Value{Omitempty: false}
	a.False(isOmitempty(v))

	v = node.NewValue("elem", reflect.ValueOf(int(0)), true, "usage")
	a.True(isOmitempty(v))
	v.Value = reflect.ValueOf(int(5))
	a.False(isOmitempty(v))

	v.Value = reflect.ValueOf(uint(0))
	a.True(isOmitempty(v))
	v.Value = reflect.ValueOf(uint(5))
	a.False(isOmitempty(v))

	v.Value = reflect.ValueOf(float64(0))
	a.True(isOmitempty(v))
	v.Value = reflect.ValueOf(float32(5))
	a.False(isOmitempty(v))

	v.Value = reflect.ValueOf([]byte{})
	a.True(isOmitempty(v))
	v.Value = reflect.ValueOf([]byte{0})
	a.False(isOmitempty(v))

	v.Value = reflect.ValueOf(false)
	a.True(isOmitempty(v))
	v.Value = reflect.ValueOf(true)
	a.False(isOmitempty(v))

	v.Value = reflect.ValueOf(map[string]string{})
	a.True(isOmitempty(v))
	v.Value = reflect.ValueOf(map[string]string{"id": "0"})
	a.False(isOmitempty(v))
}
