// SPDX-License-Identifier: MIT

package doc

import (
	"encoding/xml"
	"testing"

	"github.com/issue9/assert"
)

var _ xml.Unmarshaler = &Param{}

func TestParam_UnmarshalXML(t *testing.T) {
	a := assert.New(t)

	obj := &Param{
		Name: "text",
		Type: String,
	}
	str := `<Param name="text" type="string"></Param>`

	data, err := xml.Marshal(obj)
	a.NotError(err).
		Equal(string(data), str).
		False(obj.Optional)

	obj1 := &Param{}
	a.NotError(xml.Unmarshal([]byte(str), obj1))
	a.Equal(obj1, obj)

	// 正常
	str = `<Param name="user" deprecated="1.1.1" type="object" array="true">
		<param name="name" type="string" />
		<param name="sex" type="string">
			<enum value="male">Male</enum>
			<enum value="female">Female</enum>
		</param>
		<param name="age" type="number" />
	</Param>`
	a.NotError(xml.Unmarshal([]byte(str), obj1)).
		True(obj1.Array).
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

	// type=object，且没有子项
	obj1 = &Param{}
	str = `<Param name="v1" type="Object"></Param>`
	a.Error(xml.Unmarshal([]byte(str), obj1))

	// 相同的子项
	obj1 = &Param{}
	str = `<Param name="v1" type="Object">
		<param name="n1" type="string" />
		<param name="n1" type="number" />
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
			<enum value="female">Male</enum>
			<enum value="female">Female</enum>
	</Param>`
	a.Error(xml.Unmarshal([]byte(str), obj))
}
