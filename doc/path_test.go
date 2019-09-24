// SPDX-License-Identifier: MIT

package doc

import (
	"encoding/xml"
	"testing"

	"github.com/issue9/assert"
)

func TestPath_UnmarshalXML(t *testing.T) {
	a := assert.New(t)

	obj := &Path{
		Path:      "/users/{id}",
		Reference: "#get-users",
		Params:    []*Param{&Param{Name: "id", Type: Number}},
	}
	str := `<Path path="/users/{id}" ref="#get-users"><param name="id" type="number"></param></Path>`

	data, err := xml.Marshal(obj)
	a.NotError(err).
		Equal(string(data), str).
		Equal(obj.Path, "/users/{id}").
		Equal(obj.Reference, "#get-users").
		Equal(len(obj.Params), 1).
		Equal(obj.Params[0].Name, "id")

	obj1 := &Path{}
	a.NotError(xml.Unmarshal([]byte(str), obj1))
	a.Equal(obj1, obj)

	// query
	str = `<Path path="/users/{id}">
		<param name="id" type="number" />
		<query name="text" type="string" />
		<query name="sex" type="string">
			<enum value="male">male</enum>
			<enum value="female">female</enum>
		</query>
	</Path>`
	a.NotError(xml.Unmarshal([]byte(str), obj1)).
		Equal(len(obj1.Params), 1).
		Equal(len(obj1.Queries), 2).
		True(obj1.Queries[1].IsEnum()).
		Equal(obj1.Queries[1].Name, "sex")

	// 少 param
	str = `<Path path="/users/{id}" ref="#get-users"></Path>`
	a.Error(xml.Unmarshal([]byte(str), obj1))

	// 名称不匹配
	str = `<Path path="/users/{id}">
		<param name="not-exists" type="number" />
	</Path>`
	a.Error(xml.Unmarshal([]byte(str), obj1))

	// 少 path
	str = `<Path url="url">desc</Path>`
	a.Error(xml.Unmarshal([]byte(str), obj1))

	// 相同的参数名
	str = `<Path path="/users/{id}/logs/{id}">
		<param name="id" type="number" />
		<param name="id" type="number" />
	</Path>`
	a.Error(xml.Unmarshal([]byte(str), obj1))

	// path 语法错误
	str = `<Path path="/users/{id"></Path>`
	a.Error(xml.Unmarshal([]byte(str), obj1))

	// 语法错误
	str = `<Path path="/users/{id}">/Path>`
	a.Error(xml.Unmarshal([]byte(str), obj1))
}

func TestParsePath(t *testing.T) {
	a := assert.New(t)

	data := []*struct {
		input  string
		result []string
		err    bool
	}{
		{
			input:  "",
			result: nil,
		},
		{
			input:  "/path",
			result: nil,
		},
		{
			input:  "/users/{id}",
			result: []string{"id"},
		},
		{
			input:  "/users/{id}/logs/{lid}",
			result: []string{"id", "lid"},
		},

		{
			input: "/users/{{id}/logs/{lid}",
			err:   true,
		},
		{
			input: "/users/{id}/logs}/{lid}",
			err:   true,
		},
		{
			input: "/users/{id}/logs/{lid",
			err:   true,
		},
		{
			input: "/users/id}/logs/{lid}",
			err:   true,
		},
	}

	for index, item := range data {
		params, err := parsePath(item.input)
		if item.err {
			a.Error(err).Nil(params)
			continue
		}

		for _, param := range item.result {
			_, found := params[param]
			a.True(found, "not found @%d,v1=%s,v2=%s", index, param)
		}
	}
}
