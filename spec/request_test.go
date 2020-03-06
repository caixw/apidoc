// SPDX-License-Identifier: MIT

package spec

import (
	"encoding/xml"
	"testing"

	"github.com/issue9/assert"
)

var _ xml.Unmarshaler = &Request{}

func TestRequest_UnmarshalXML(t *testing.T) {
	a := assert.New(t)

	obj := &Request{
		Type:     String,
		Mimetype: "json",
	}
	str := `<Request type="string" mimetype="json"></Request>`

	data, err := xml.Marshal(obj)
	a.NotError(err).Equal(string(data), str)

	obj1 := &Request{}
	a.NotError(xml.Unmarshal([]byte(str), obj1))
	a.Equal(obj1, obj)

	// 正常
	obj1 = &Request{}
	str = `<Request deprecated="1.1.1" type="object" array="true" mimetype="json">
		<param name="name" type="string" summary="name" />
		<param name="sex" type="string" summary="sex">
			<enum value="male" summary="male" />
			<enum value="female" summary="female" />
		</param>
		<param name="age" type="number" summary="age" />
	</Request>`
	a.NotError(xml.Unmarshal([]byte(str), obj1)).
		True(obj1.Array).
		Equal(obj1.Type, Object).
		Equal(obj1.Deprecated, "1.1.1").
		Equal(3, len(obj1.Items))

	obj1 = &Request{}
	str = `<Request type="string"></Request>`
	a.NotError(xml.Unmarshal([]byte(str), obj1))

	// type=object，且没有子项
	obj1 = &Request{}
	str = `<Request type="Object" mimetype="json"></Request>`
	a.Error(xml.Unmarshal([]byte(str), obj1))

	// 相同的子项
	obj1 = &Request{}
	str = `<Request type="Object" mimetype="json">
		<param name="n1" type="string" summary="n1" />
		<param name="n1" type="number" summary="n2" />
	</Request>`
	a.Error(xml.Unmarshal([]byte(str), obj1))

	// 报头不能为 Object
	obj1 = &Request{}
	str = `<Request type="number" mimetype="json">
		<header name="k1" type="string" summary="n1" />
		<header name="k2" type="object" summary="n2">
			<param name="xx" type="number" summary="summary" />
		</header>
	</Request>`
	a.Error(xml.Unmarshal([]byte(str), obj1))

	// example 不匹配
	obj1 = &Request{}
	str = `<Request type="Object" mimetype="json">
		<param name="n1" type="string" summary="n1" />
		<example mimetype="xml"><![CDATA[xx]]></example>
	</Request>`
	a.Error(xml.Unmarshal([]byte(str), obj1))

	// 语法错误
	obj1 = &Request{}
	str = `<Request deprecated="x.1.1" mimetype="json">text</Request>`
	a.Error(xml.Unmarshal([]byte(str), obj1))
}

func TestRequest_UnmarshalXML_enum(t *testing.T) {
	a := assert.New(t)

	obj := &Request{}
	str := `<Request name="sex" type="string" mimetype="json">
			<enum value="male" summary="male" />
			<enum value="female" summary="female" />
	</Request>`
	a.NotError(xml.Unmarshal([]byte(str), obj)).
		False(obj.Array).
		True(obj.IsEnum()).
		Equal(obj.Type, String).
		Equal(2, len(obj.Enums))

	// 枚举中存在相同值
	obj = &Request{}
	str = `<Request name="sex" type="string" mimetype="json">
			<enum value="female" summary="male" />
			<enum value="female" summary="female" />
	</Request>`
	a.Error(xml.Unmarshal([]byte(str), obj))
}
