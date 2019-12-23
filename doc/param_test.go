// SPDX-License-Identifier: MIT

package doc

import (
	"encoding/xml"
	"testing"

	"github.com/issue9/assert"
)

var (
	_ xml.Unmarshaler = &Param{}
)

func TestParam_UnmarshalXML(t *testing.T) {
	a := assert.New(t)

	obj := &Param{
		Name:    "text",
		Type:    String,
		Summary: "text",
		XML:     XML{XMLAttr: true},
	}
	str := `<Param xml-attr="true" name="text" type="string" summary="text"></Param>`

	data, err := xml.Marshal(obj)
	a.NotError(err).
		Equal(string(data), str).
		False(obj.Optional)

	obj1 := &Param{}
	a.NotError(xml.Unmarshal([]byte(str), obj1))
	a.Equal(obj1, obj)

	// 正常
	str = `<Param name="user" deprecated="1.1.1" type="object" array="true">
		<description><![CDATA[<a>user</a>]]></description>
		<param name="name" type="string" summary="name" />
		<param name="sex" type="string" summary="sex">
			<enum value="male" summary="Male" />
			<enum value="female" summary="female" />
		</param>
		<param name="age" type="number" summary="age" />
	</Param>`
	obj1 = &Param{}
	a.NotError(xml.Unmarshal([]byte(str), obj1)).
		True(obj1.Array).
		Equal(obj1.Description.Text, "<a>user</a>").
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

func TestParam_UnmarshalXML_enum(t *testing.T) {
	a := assert.New(t)

	obj := &Param{}
	str := `<Param name="sex" type="string">
			<description>sex</description>
			<enum value="male" summary="Male" />
			<enum value="female" summary="female" />
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
			<enum value="male" summary="Male" />
			<enum value="male" summary="female" />
	</Param>`
	a.Error(xml.Unmarshal([]byte(str), obj))

	// 枚举类型与父类型不相同
	obj = &Param{}
	str = `<Param name="sex" type="number">
			<description>sex</description>
			<enum value="1" summary="Male" />
			<enum value="male" summary="female" />
	</Param>`
	a.Error(xml.Unmarshal([]byte(str), obj))
}

func TestChkEnumsType(t *testing.T) {
	a := assert.New(t)

	boolEnums := []*Enum{
		{Value: "true"},
		{Value: "false"},
	}
	numberEnums := []*Enum{
		{Value: "1"},
		{Value: "2"},
	}
	stringEnums := []*Enum{
		{Value: "1"},
		{Value: "true"},
		{Value: "string"},
	}

	a.NotError(chkEnumsType(Bool, boolEnums, ""))
	a.Error(chkEnumsType(Bool, numberEnums, ""))
	a.Error(chkEnumsType(Bool, stringEnums, ""))

	a.NotError(chkEnumsType(Number, numberEnums, ""))
	a.Error(chkEnumsType(Number, boolEnums, ""))
	a.Error(chkEnumsType(Number, stringEnums, ""))

	a.NotError(chkEnumsType(String, numberEnums, ""))
	a.NotError(chkEnumsType(String, boolEnums, ""))
	a.NotError(chkEnumsType(String, stringEnums, ""))

	a.Error(chkEnumsType(None, numberEnums, ""))
	a.Error(chkEnumsType(None, boolEnums, ""))
	a.Error(chkEnumsType(None, stringEnums, ""))

	a.Error(chkEnumsType(Object, numberEnums, ""))
	a.Error(chkEnumsType(Object, boolEnums, ""))
	a.Error(chkEnumsType(Object, stringEnums, ""))
}
