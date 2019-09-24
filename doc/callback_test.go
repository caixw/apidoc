// SPDX-License-Identifier: MIT

package doc

import (
	"encoding/xml"
	"net/http"
	"testing"

	"github.com/issue9/assert"
)

var _ xml.Unmarshaler = &Callback{}

func TestCallback_UnmarshalXML(t *testing.T) {
	a := assert.New(t)

	obj := &Callback{
		Schema:   "http",
		Method:   http.MethodGet,
		Requests: []*Request{&Request{Mimetype: "json", Type: String}},
	}
	str := `<Callback schema="http" method="GET"><request type="string" mimetype="json"></request></Callback>`

	data, err := xml.Marshal(obj)
	a.NotError(err).Equal(string(data), str)

	obj1 := &Callback{}
	a.NotError(xml.Unmarshal([]byte(str), obj1))
	a.Equal(obj1, obj)

	// 正常
	str = `<Callback deprecated="1.1.1" method="GET" schema="HTTPS">
		<request status="200" mimetype="json" type="object">
			<param name="name" type="string" />
			<param name="sex" type="string">
				<enum value="male">Male</enum>
				<enum value="female">Female</enum>
			</param>
			<param name="age" type="number" />
		</request>
	</Callback>`
	a.NotError(xml.Unmarshal([]byte(str), obj1)).
		Equal(obj1.Deprecated, "1.1.1").
		Equal(1, len(obj1.Requests)).
		Equal(obj1.Requests[0].Type, Object)

	// 少 schema
	str = `<Callback method="GET"><request type="string" mimetype="json" /></Callback>`
	a.Error(xml.Unmarshal([]byte(str), obj1))

	// 少 method
	str = `<Callback schema="HTTP"><request type="string" mimetype="json" /></Callback>`
	a.Error(xml.Unmarshal([]byte(str), obj1))

	// 无效的 schema
	str = `<Callback method="GET" schema="invalid"><request type="string" mimetype="json" /></Callback>`
	a.Error(xml.Unmarshal([]byte(str), obj1))

	// 缺少 request
	str = `<Callback method="GET" schema="http"></Callback>`
	a.Error(xml.Unmarshal([]byte(str), obj1))

	// 语法错误
	str = `<Callback name="url" deprecated="x.1.1">text</Callback>`
	a.Error(xml.Unmarshal([]byte(str), obj1))
}
