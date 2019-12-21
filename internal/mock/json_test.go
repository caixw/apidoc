// SPDX-License-Identifier: MIT

package mock

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v5/doc"
)

type jsonTester struct {
	Title string
	Type  *doc.Request
	Data  string
	Err   bool
}

var jsonTestData = []*jsonTester{
	{
		Title: "nil",
		Type:  nil,
		Data:  "null",
	},
	{
		Title: "doc.None",
		Type:  &doc.Request{Type: doc.None},
		Data:  "",
	},
	{
		Title: "doc.Request{}",
		Type:  &doc.Request{},
		Data:  "",
	},
	{
		Title: "number",
		Type:  &doc.Request{Type: doc.Number},
		Data:  "1024",
	},
	{
		Title: "enum number",
		Type: &doc.Request{
			Type: doc.Number,
			Enums: []*doc.Enum{
				{Value: "1024"},
				{Value: "1025"},
			},
		},
		Data: "1024",
	},
	{
		Title: "enum string",
		Type: &doc.Request{
			Type: doc.String,
			Enums: []*doc.Enum{
				{Value: "1024"},
				{Value: "1025"},
			},
		},
		Data: `"1024"`,
	},
	{ // array
		Title: "[bool]",
		Type: &doc.Request{
			Type:  doc.Bool,
			Array: true,
		},
		Data: `[
    true,
    true,
    true,
    true,
    true
]`,
	},
	{
		Title: "array with enum",
		Type: &doc.Request{
			Type:  doc.Number,
			Array: true,
			Enums: []*doc.Enum{
				{Value: "1"},
				{Value: "2"},
				{Value: "3"},
			},
		},
		Data: `[
    1,
    1,
    1,
    1,
    1
]`,
	},
	{
		Title: "bool",
		Type:  &doc.Request{Type: doc.Bool},
		Data:  "true",
	},
	{ // Object
		Title: "Object",
		Type: &doc.Request{
			Type: doc.Object,
			Items: []*doc.Param{
				{
					Type: doc.String,
					Name: "name",
				},
				{
					Type: doc.Number,
					Name: "id",
				},
			},
		},
		Data: `{
    "name": "1024",
    "id": 1024
}`,
	},

	{ // 各类型混合
		Title: "Object with array",
		Type: &doc.Request{
			Type: doc.Object,
			Items: []*doc.Param{
				{
					Type: doc.String,
					Name: "name",
				},
				{
					Type: doc.Number,
					Name: "id",
				},
				{
					Type: doc.Object,
					Name: "group",
					Items: []*doc.Param{
						{
							Type: doc.String,
							Name: "name",
						},
						{
							Type: doc.Number,
							Name: "id",
						},
						{
							Name:  "tags",
							Array: true,
							Type:  doc.Object,
							Items: []*doc.Param{
								{
									Type: doc.String,
									Name: "name",
								},
								{
									Type: doc.Number,
									Name: "id",
								},
							},
						}, // end tags
					},
				}, // end group
			},
		},
		Data: `{
    "name": "1024",
    "id": 1024,
    "group": {
        "name": "1024",
        "id": 1024,
        "tags": [
            {
                "name": "1024",
                "id": 1024
            },
            {
                "name": "1024",
                "id": 1024
            },
            {
                "name": "1024",
                "id": 1024
            },
            {
                "name": "1024",
                "id": 1024
            },
            {
                "name": "1024",
                "id": 1024
            }
        ]
    }
}`,
	},
}

func TestJSONValidator_find(t *testing.T) {
	item := jsonTestData[len(jsonTestData)-1]

	a := assert.New(t)
	v := &jsonValidator{
		param:   item.Type.ToParam(),
		decoder: json.NewDecoder(strings.NewReader(item.Data)),
	}

	v.names = []string{}
	p := v.find()
	a.Equal(p, v.param)

	v.names = nil
	p = v.find()
	a.Equal(p, v.param)

	v.names = []string{""}
	p = v.find()
	a.Nil(p)

	v.names = []string{"name"}
	p = v.find()
	a.NotNil(p).Equal(p.Type, doc.String)

	v.names = []string{"not-exists"}
	p = v.find()
	a.Nil(p)

	v.names = []string{"group", "id"}
	p = v.find()
	a.NotNil(p).Equal(p.Type, doc.Number)

	v.names = []string{"group", "tags", "id"}
	p = v.find()
	a.NotNil(p).Equal(p.Type, doc.Number)
}

func TestValidJSON(t *testing.T) {
	a := assert.New(t)

	for _, item := range jsonTestData {
		err := validJSON(item.Type, []byte(item.Data))
		if item.Err {
			a.Error(err)
		} else {
			a.NotError(err, "测试 %s 时返回错误值 %s", item.Title, err)
		}
	}
}

func TestBuildJSON(t *testing.T) {
	a := assert.New(t)

	for _, item := range jsonTestData {
		data, err := buildJSON(item.Type)

		if item.Err {
			a.Error(err).Nil(data)
		} else {
			a.NotError(err, "测试 %s 返回了错误值 %s", item.Title, err).
				Equal(string(data), item.Data, "测试 %s 失败 v1:%s,v2:%s", item.Title, string(data), item.Data)
		}
	}
}
