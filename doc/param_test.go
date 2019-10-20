// SPDX-License-Identifier: MIT

package doc

import (
	"encoding/xml"
	"testing"

	"github.com/issue9/assert"
)

var (
	_ xml.Unmarshaler = &Param{}
	_ xml.Unmarshaler = &SimpleParam{}
)

func TestParam_UnmarshalXML(t *testing.T) {
	a := assert.New(t)

	obj := &Param{
		Name:    "text",
		Type:    String,
		Summary: "text",
	}
	str := `<Param name="text" type="string" summary="text"></Param>`

	data, err := xml.Marshal(obj)
	a.NotError(err).
		Equal(string(data), str).
		False(obj.Optional)

	obj1 := &Param{}
	a.NotError(xml.Unmarshal([]byte(str), obj1))
	a.Equal(obj1, obj)

	// 正常
	str = `<Param name="user" deprecated="1.1.1" type="object" array="true">
		<description><![CDATA[user]]></description>
		<param name="name" type="string" summary="name" />
		<param name="sex" type="string" summary="sex">
			<enum value="male">Male</enum>
			<enum value="female">Female</enum>
		</param>
		<param name="age" type="number" summary="age" />
	</Param>`
	a.NotError(xml.Unmarshal([]byte(str), obj1)).
		True(obj1.Array).
		Equal(obj1.Description.String(), "user").
		Equal(obj1.Type, Object).
		Equal(obj1.Deprecated, "1.1.1").
		Equal(3, len(obj1.Items)).
		False(obj1.Optional)

	// 少 name
	obj1 = &Param{}
	str = `<Param url="url">desc</Param>`
	a.Error(xml.Unmarshal([]byte(str), obj1))

	// 少 type
	obj1 = &Param{}
	str = `<Param name="v1"></Param>`
	a.Error(xml.Unmarshal([]byte(str), obj1))

	// 少 summary 和 description
	obj1 = &Param{}
	str = `<Param name="v1" type="string"></Param>`
	a.Error(xml.Unmarshal([]byte(str), obj1))

	// type=object，且没有子项
	obj1 = &Param{}
	str = `<Param name="v1" type="Object"></Param>`
	a.Error(xml.Unmarshal([]byte(str), obj1))

	// 相同的子项
	obj1 = &Param{}
	str = `<Param name="v1" type="Object">
		<param name="n1" type="string" summary="n1" />
		<param name="n1" type="number" summary="n2" />
	</Param>`
	a.Error(xml.Unmarshal([]byte(str), obj1))

	// 语法错误
	obj1 = &Param{}
	str = `<Param name="url" deprecated="x.1.1">text</Param>`
	a.Error(xml.Unmarshal([]byte(str), obj1))
}

func TestSimpleParam_UnmarshalXML(t *testing.T) {
	a := assert.New(t)

	obj := &SimpleParam{
		Name:        "text",
		Type:        String,
		Description: "test",
	}
	str := `<SimpleParam name="text" type="string"><![CDATA[test]]></SimpleParam>`

	data, err := xml.Marshal(obj)
	a.NotError(err).
		Equal(string(data), str).
		False(obj.Optional)

	obj1 := &SimpleParam{}
	a.NotError(xml.Unmarshal([]byte(str), obj1))
	a.Equal(obj1, obj)

	// 正常
	str = `<SimpleParam name="user" deprecated="1.1.1" type="string" summary="user"></SimpleParam>`
	a.NotError(xml.Unmarshal([]byte(str), obj1)).
		Equal(obj1.Type, String).
		Equal(obj1.Summary, "user").
		Equal(obj1.Deprecated, "1.1.1").
		False(obj1.Optional)

	// 少 name
	obj1 = &SimpleParam{}
	str = `<SimpleParam url="url">desc</SimpleParam>`
	a.Error(xml.Unmarshal([]byte(str), obj1))

	// 少 type
	obj1 = &SimpleParam{}
	str = `<SimpleParam name="v1"></SimpleParam>`
	a.Error(xml.Unmarshal([]byte(str), obj1))

	// type=object
	obj1 = &SimpleParam{}
	str = `<SimpleParam name="v1" type="Object"></SimpleParam>`
	a.Error(xml.Unmarshal([]byte(str), obj1))

	// 语法错误
	obj1 = &SimpleParam{}
	str = `<SimpleParam name="url" deprecated="x.1.1">text</SimpleParam>`
	a.Error(xml.Unmarshal([]byte(str), obj1))
}

func TestParam_UnmarshalXML_enum(t *testing.T) {
	a := assert.New(t)

	obj := &Param{}
	str := `<Param name="sex" type="string">
			<description>sex</description>
			<enum value="male">Male</enum>
			<enum value="female">Female</enum>
	</Param>`
	a.NotError(xml.Unmarshal([]byte(str), obj)).
		False(obj.Array).
		True(obj.IsEnum()).
		Equal(obj.Type, String).
		Equal(2, len(obj.Enums))

	// 枚举中存在相同值
	obj = &Param{}
	str = `<Param name="sex" type="string">
			<description>sex</description>
			<enum value="female">Male</enum>
			<enum value="female">Female</enum>
	</Param>`
	a.Error(xml.Unmarshal([]byte(str), obj))
}

func TestSimpleParam_UnmarshalXML_enum(t *testing.T) {
	a := assert.New(t)

	obj := &SimpleParam{}
	str := `<SimpleParam name="sex" type="string">
			<description>sex</description>
			<enum value="male">Male</enum>
			<enum value="female">Female</enum>
	</SimpleParam>`
	a.NotError(xml.Unmarshal([]byte(str), obj)).
		False(obj.Array).
		True(obj.IsEnum()).
		Equal(obj.Type, String).
		Equal(2, len(obj.Enums))

	// 枚举中存在相同值
	obj = &SimpleParam{}
	str = `<SimpleParam name="sex" type="string">
			<description>sex</description>
			<enum value="female">Male</enum>
			<enum value="female">Female</enum>
	</SimpleParam>`
	a.Error(xml.Unmarshal([]byte(str), obj))
}
