// SPDX-License-Identifier: MIT

package mock

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/issue9/assert"
	"github.com/issue9/assert/rest"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/core/messagetest"
	"github.com/caixw/apidoc/v7/internal/ast"
	"github.com/caixw/apidoc/v7/internal/ast/asttest"
	"github.com/caixw/apidoc/v7/internal/xmlenc"
)

var _ http.Handler = &mock{}

type tester struct {
	Title string
	Type  *ast.Request
	JSON  string
	XML   string
	XMLNS []*ast.XMLNamespace
}

var dataWithHeader = &tester{
	Title: "object with item",
	Type: &ast.Request{
		Name: &ast.Attribute{Value: xmlenc.String{Value: "root"}},
		Type: &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeObject}},
		Headers: []*ast.Param{
			{
				Type: &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeString}},
				Name: &ast.Attribute{Value: xmlenc.String{Value: "content-type"}},
			},
			{
				Type: &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeString}},
				Name: &ast.Attribute{Value: xmlenc.String{Value: "encoding"}},
			},
		},
		Items: []*ast.Param{
			{
				Type: &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeObject}},
				Name: &ast.Attribute{Value: xmlenc.String{Value: "name"}},
				Items: []*ast.Param{
					{
						Type: &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeString}},
						Name: &ast.Attribute{Value: xmlenc.String{Value: "last"}},
					},
					{
						Type:     &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeString}},
						Name:     &ast.Attribute{Value: xmlenc.String{Value: "first"}},
						Optional: &ast.BoolAttribute{Value: ast.Bool{Value: true}},
					},
				},
			},
			{
				Type: &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeNumber}},
				Name: &ast.Attribute{Value: xmlenc.String{Value: "age"}},
				XML:  ast.XML{XMLAttr: &ast.BoolAttribute{Value: ast.Bool{Value: true}}},
			},
		},
	},
	JSON: `{
    "name": {
        "last": "1024",
        "first": "1024"
    },
    "age": 1024
}`,
	XML: `<root age="1024">
    <name>
        <last>1024</last>
        <first>1024</first>
    </name>
</root>`,
}

// 提供了测试 validJSON/buildXML 和 buildJSON/buildXML 的数据
var data = []*tester{
	{
		Title: "nil",
		Type:  nil,
		JSON:  "null",
		XML:   "",
	},
	{
		Title: "doc.None",
		Type:  &ast.Request{Type: &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeNone}}},
		JSON:  "",
		XML:   "",
	},
	{
		Title: "doc.Request{}",
		Type:  &ast.Request{},
		JSON:  "",
		XML:   "",
	},
	{
		Title: "number",
		Type: &ast.Request{
			Type: &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeNumber}},
			Name: &ast.Attribute{Value: xmlenc.String{Value: "root"}},
		},
		JSON: "1024",
		XML:  "<root>1024</root>",
	},
	{
		Title: "number-cdata",
		Type: &ast.Request{
			Type: &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeNumber}},
			Name: &ast.Attribute{Value: xmlenc.String{Value: "root"}},
			XML:  ast.XML{XMLCData: &ast.BoolAttribute{Value: ast.Bool{Value: true}}},
		},
		JSON: "1024",
		XML:  "<root><![CDATA[1024]]></root>",
	},
	{
		Title: "date-cdata",
		Type: &ast.Request{
			Type: &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeDate}},
			Name: &ast.Attribute{Value: xmlenc.String{Value: "root"}},
			XML:  ast.XML{XMLCData: &ast.BoolAttribute{Value: ast.Bool{Value: true}}},
		},
		JSON: `"2020-01-02"`,
		XML:  "<root><![CDATA[2020-01-02]]></root>",
	},
	{
		Title: "time-cdata",
		Type: &ast.Request{
			Type: &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeTime}},
			Name: &ast.Attribute{Value: xmlenc.String{Value: "root"}},
			XML:  ast.XML{XMLCData: &ast.BoolAttribute{Value: ast.Bool{Value: true}}},
		},
		JSON: `"15:16:17Z"`,
		XML:  "<root><![CDATA[15:16:17Z]]></root>",
	},
	{
		Title: "date-time-cdata",
		Type: &ast.Request{
			Type: &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeDateTime}},
			Name: &ast.Attribute{Value: xmlenc.String{Value: "root"}},
			XML:  ast.XML{XMLCData: &ast.BoolAttribute{Value: ast.Bool{Value: true}}},
		},
		JSON: `"2020-01-02T15:16:17Z"`,
		XML:  "<root><![CDATA[2020-01-02T15:16:17Z]]></root>",
	},
	{ // namespace
		Title: "number-cdata-namespace",
		Type: &ast.Request{
			Type: &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeNumber}},
			Name: &ast.Attribute{Value: xmlenc.String{Value: "root"}},
			XML: ast.XML{
				XMLCData:    &ast.BoolAttribute{Value: ast.Bool{Value: true}},
				XMLNSPrefix: &ast.Attribute{Value: xmlenc.String{Value: "ns"}},
			},
		},
		JSON: "1024",
		XML:  `<ns:root xmlns:ns="urn"><![CDATA[1024]]></ns:root>`,
		XMLNS: []*ast.XMLNamespace{
			{
				Prefix: &ast.Attribute{Value: xmlenc.String{Value: "ns"}},
				URN:    &ast.Attribute{Value: xmlenc.String{Value: "urn"}},
			},
		},
	},
	{
		Title: "enum number",
		Type: &ast.Request{
			Name: &ast.Attribute{Value: xmlenc.String{Value: "root"}},
			Type: &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeNumber}},
			Enums: []*ast.Enum{
				{Value: &ast.Attribute{Value: xmlenc.String{Value: "1024"}}},
				{Value: &ast.Attribute{Value: xmlenc.String{Value: "1025"}}},
			},
		},
		JSON: "1024",
		XML:  `<root>1024</root>`,
	},
	{
		Title: "xml-extract",
		Type: &ast.Request{
			Name: &ast.Attribute{Value: xmlenc.String{Value: "root"}},
			Type: &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeObject}},
			Items: []*ast.Param{
				{
					Name: &ast.Attribute{Value: xmlenc.String{Value: "id"}},
					Type: &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeInt}},
					XML:  ast.XML{XMLAttr: &ast.BoolAttribute{Value: ast.Bool{Value: true}}},
				},
				{
					Name: &ast.Attribute{Value: xmlenc.String{Value: "desc"}},
					Type: &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeString}},
					XML:  ast.XML{XMLExtract: &ast.BoolAttribute{Value: ast.Bool{Value: true}}},
				},
			},
		},
		JSON: `{
    "id": 1024,
    "desc": "1024"
}`,
		XML: `<root id="1024">1024</root>`,
	},

	{
		Title: "xml-extract-cdata",
		Type: &ast.Request{
			Name: &ast.Attribute{Value: xmlenc.String{Value: "root"}},
			Type: &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeObject}},
			Items: []*ast.Param{
				{
					Name: &ast.Attribute{Value: xmlenc.String{Value: "id"}},
					Type: &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeNumber}},
					XML:  ast.XML{XMLAttr: &ast.BoolAttribute{Value: ast.Bool{Value: true}}},
				},
				{
					Name: &ast.Attribute{Value: xmlenc.String{Value: "desc"}},
					Type: &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeString}},
					XML: ast.XML{
						XMLExtract: &ast.BoolAttribute{Value: ast.Bool{Value: true}},
						XMLCData:   &ast.BoolAttribute{Value: ast.Bool{Value: true}},
					},
				},
			},
		},
		JSON: `{
    "id": 1024,
    "desc": "1024"
}`,
		XML: `<root id="1024"><![CDATA[1024]]></root>`,
	},

	{
		Title: "enum string",
		Type: &ast.Request{
			Name: &ast.Attribute{Value: xmlenc.String{Value: "root"}},
			Type: &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeString}},
			Enums: []*ast.Enum{
				{Value: &ast.Attribute{Value: xmlenc.String{Value: "1024"}}},
				{Value: &ast.Attribute{Value: xmlenc.String{Value: "1025"}}},
			},
		},
		JSON: `"1024"`,
		XML:  `<root>1024</root>`,
	},
	{ // array
		Title: "[bool] with wrapped",
		Type: &ast.Request{
			XML:   ast.XML{XMLWrapped: &ast.Attribute{Value: xmlenc.String{Value: "root"}}},
			Name:  &ast.Attribute{Value: xmlenc.String{Value: "arr"}},
			Type:  &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeBool}},
			Array: &ast.BoolAttribute{Value: ast.Bool{Value: true}},
		},
		JSON: `[
    true,
    true,
    true,
    true,
    true
]`,
		XML: `<root>
    <arr>true</arr>
    <arr>true</arr>
    <arr>true</arr>
    <arr>true</arr>
    <arr>true</arr>
</root>`,
	},
	{
		Title: "[bool] with wrapped>",
		Type: &ast.Request{
			XML: ast.XML{
				XMLWrapped:  &ast.Attribute{Value: xmlenc.String{Value: "root>array"}},
				XMLNSPrefix: &ast.Attribute{Value: xmlenc.String{Value: "ns"}},
			},
			Name:  &ast.Attribute{Value: xmlenc.String{Value: "arr"}},
			Type:  &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeBool}},
			Array: &ast.BoolAttribute{Value: ast.Bool{Value: true}},
		},
		XMLNS: []*ast.XMLNamespace{
			{
				URN:    &ast.Attribute{Value: xmlenc.String{Value: "urn"}},
				Prefix: &ast.Attribute{Value: xmlenc.String{Value: "ns"}},
			},
		},
		JSON: `[
    true,
    true,
    true,
    true,
    true
]`,
		XML: `<ns:root xmlns:ns="urn">
    <ns:array>true</ns:array>
    <ns:array>true</ns:array>
    <ns:array>true</ns:array>
    <ns:array>true</ns:array>
    <ns:array>true</ns:array>
</ns:root>`,
	},
	{
		Title: "[bool] with >wrapped",
		Type: &ast.Request{
			Name: &ast.Attribute{Value: xmlenc.String{Value: "bool"}},
			Type: &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeObject}},
			Items: []*ast.Param{
				{
					XML:   ast.XML{XMLWrapped: &ast.Attribute{Value: xmlenc.String{Value: ">array"}}},
					Name:  &ast.Attribute{Value: xmlenc.String{Value: "arr"}},
					Type:  &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeBool}},
					Array: &ast.BoolAttribute{Value: ast.Bool{Value: true}},
				},
			},
		},
		JSON: `{
    "arr": [
        true,
        true,
        true,
        true,
        true
    ]
}`,
		XML: `<bool>
    <array>true</array>
    <array>true</array>
    <array>true</array>
    <array>true</array>
    <array>true</array>
</bool>`,
	},
	{
		Title: "array with enum",
		Type: &ast.Request{
			XML:   ast.XML{XMLWrapped: &ast.Attribute{Value: xmlenc.String{Value: "root"}}},
			Name:  &ast.Attribute{Value: xmlenc.String{Value: "arr"}},
			Type:  &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeNumber}},
			Array: &ast.BoolAttribute{Value: ast.Bool{Value: true}},
			Enums: []*ast.Enum{
				{Value: &ast.Attribute{Value: xmlenc.String{Value: "1"}}},
				{Value: &ast.Attribute{Value: xmlenc.String{Value: "2"}}},
				{Value: &ast.Attribute{Value: xmlenc.String{Value: "3"}}},
			},
		},
		JSON: `[
    1,
    1,
    1,
    1,
    1
]`,
		XML: `<root>
    <arr>1</arr>
    <arr>1</arr>
    <arr>1</arr>
    <arr>1</arr>
    <arr>1</arr>
</root>`,
	},
	{
		Title: "bool",
		Type: &ast.Request{
			Type: &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeBool}},
			Name: &ast.Attribute{Value: xmlenc.String{Value: "root"}},
		},
		JSON: "true",
		XML:  "<root>true</root>",
	},
	{ // Object
		Title: "Object with wrapped",
		Type: &ast.Request{
			Name: &ast.Attribute{Value: xmlenc.String{Value: "root"}},
			Type: &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeObject}},
			Items: []*ast.Param{
				{
					Type: &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeString}},
					Name: &ast.Attribute{Value: xmlenc.String{Value: "name"}},
				},
				{
					Type:  &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeNumber}},
					Name:  &ast.Attribute{Value: xmlenc.String{Value: "num"}},
					Array: &ast.BoolAttribute{Value: ast.Bool{Value: true}},
					XML:   ast.XML{XMLWrapped: &ast.Attribute{Value: xmlenc.String{Value: "nums"}}},
				},
				{
					Type: &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeNumber}},
					Name: &ast.Attribute{Value: xmlenc.String{Value: "id"}},
					XML:  ast.XML{XMLAttr: &ast.BoolAttribute{Value: ast.Bool{Value: true}}},
				},
			},
		},
		JSON: `{
    "name": "1024",
    "num": [
        1024,
        1024,
        1024,
        1024,
        1024
    ],
    "id": 1024
}`,
		XML: `<root id="1024">
    <name>1024</name>
    <nums>
        <num>1024</num>
        <num>1024</num>
        <num>1024</num>
        <num>1024</num>
        <num>1024</num>
    </nums>
</root>`,
	},

	{
		Title: "object array",
		Type: &ast.Request{
			XML:   ast.XML{XMLWrapped: &ast.Attribute{Value: xmlenc.String{Value: "root"}}},
			Name:  &ast.Attribute{Value: xmlenc.String{Value: "user"}},
			Type:  &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeObject}},
			Array: &ast.BoolAttribute{Value: ast.Bool{Value: true}},
			Items: []*ast.Param{
				{
					Name: &ast.Attribute{Value: xmlenc.String{Value: "id"}},
					Type: &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeNumber}},
				},
				{
					Name: &ast.Attribute{Value: xmlenc.String{Value: "name"}},
					Type: &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeEmail}},
				},
			},
		},
		JSON: `[
    {
        "id": 1024,
        "name": "user@example.com"
    },
    {
        "id": 1024,
        "name": "user@example.com"
    },
    {
        "id": 1024,
        "name": "user@example.com"
    },
    {
        "id": 1024,
        "name": "user@example.com"
    },
    {
        "id": 1024,
        "name": "user@example.com"
    }
]`,
		XML: `<root>
    <user>
        <id>1024</id>
        <name>user@example.com</name>
    </user>
    <user>
        <id>1024</id>
        <name>user@example.com</name>
    </user>
    <user>
        <id>1024</id>
        <name>user@example.com</name>
    </user>
    <user>
        <id>1024</id>
        <name>user@example.com</name>
    </user>
    <user>
        <id>1024</id>
        <name>user@example.com</name>
    </user>
</root>`,
	},

	{
		Title: "object with item namespace",
		Type: &ast.Request{
			Name: &ast.Attribute{Value: xmlenc.String{Value: "root"}},
			Type: &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeObject}},
			Headers: []*ast.Param{
				{
					Type: &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeString}},
					Name: &ast.Attribute{Value: xmlenc.String{Value: "content-type"}},
				},
				{
					Type: &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeString}},
					Name: &ast.Attribute{Value: xmlenc.String{Value: "encoding"}},
				},
			},
			Items: []*ast.Param{
				{
					Type: &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeObject}},
					Name: &ast.Attribute{Value: xmlenc.String{Value: "name"}},
					XML:  ast.XML{XMLNSPrefix: &ast.Attribute{Value: xmlenc.String{Value: "aa"}}},
					Items: []*ast.Param{
						{
							Type: &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeString}},
							Name: &ast.Attribute{Value: xmlenc.String{Value: "last"}},
							XML:  ast.XML{XMLNSPrefix: &ast.Attribute{Value: xmlenc.String{Value: "aa"}}},
						},
						{
							Type:     &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeString}},
							Name:     &ast.Attribute{Value: xmlenc.String{Value: "first"}},
							Optional: &ast.BoolAttribute{Value: ast.Bool{Value: true}},
						},
					},
				},
				{
					Type: &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeNumber}},
					Name: &ast.Attribute{Value: xmlenc.String{Value: "age"}},
					XML: ast.XML{
						XMLAttr:     &ast.BoolAttribute{Value: ast.Bool{Value: true}},
						XMLNSPrefix: &ast.Attribute{Value: xmlenc.String{Value: "aa"}},
					},
				},
			},
		},
		XMLNS: []*ast.XMLNamespace{
			{
				URN: &ast.Attribute{Value: xmlenc.String{Value: "urn"}},
			},
			{
				Prefix: &ast.Attribute{Value: xmlenc.String{Value: "aa"}},
				URN:    &ast.Attribute{Value: xmlenc.String{Value: "aa-urn"}},
			},
		},
		JSON: `{
    "name": {
        "last": "1024",
        "first": "1024"
    },
    "age": 1024
}`,
		XML: `<root aa:age="1024" xmlns="urn" xmlns:aa="aa-urn">
    <aa:name>
        <aa:last>1024</aa:last>
        <first>1024</first>
    </aa:name>
</root>`,
	},

	{ // 各类型混合
		Title: "Object with array",
		Type: &ast.Request{
			Name: &ast.Attribute{Value: xmlenc.String{Value: "root"}},
			Type: &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeObject}},
			Items: []*ast.Param{
				{
					Type: &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeString}},
					Name: &ast.Attribute{Value: xmlenc.String{Value: "name"}},
				},
				{
					Type: &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeNumber}},
					Name: &ast.Attribute{Value: xmlenc.String{Value: "id"}},
					XML:  ast.XML{XMLAttr: &ast.BoolAttribute{Value: ast.Bool{Value: true}}},
				},
				{
					Type: &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeObject}},
					Name: &ast.Attribute{Value: xmlenc.String{Value: "group"}},
					Items: []*ast.Param{
						{
							Type: &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeString}},
							Name: &ast.Attribute{Value: xmlenc.String{Value: "name"}},
							XML:  ast.XML{XMLAttr: &ast.BoolAttribute{Value: ast.Bool{Value: true}}},
						},
						{
							Type: &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeNumber}},
							Name: &ast.Attribute{Value: xmlenc.String{Value: "id"}},
							XML:  ast.XML{XMLAttr: &ast.BoolAttribute{Value: ast.Bool{Value: true}}},
						},
						{
							Name:  &ast.Attribute{Value: xmlenc.String{Value: "tags"}},
							Array: &ast.BoolAttribute{Value: ast.Bool{Value: true}},
							Type:  &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeObject}},
							Items: []*ast.Param{
								{
									Type: &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeString}},
									Name: &ast.Attribute{Value: xmlenc.String{Value: "name"}},
								},
								{
									Type: &ast.TypeAttribute{Value: xmlenc.String{Value: ast.TypeNumber}},
									Name: &ast.Attribute{Value: xmlenc.String{Value: "id"}},
									XML:  ast.XML{XMLAttr: &ast.BoolAttribute{Value: ast.Bool{Value: true}}},
								},
							},
						}, // end tags
					},
				}, // end group
			},
		},
		JSON: `{
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
		XML: `<root id="1024">
    <name>1024</name>
    <group name="1024" id="1024">
        <tags id="1024">
            <name>1024</name>
        </tags>
        <tags id="1024">
            <name>1024</name>
        </tags>
        <tags id="1024">
            <name>1024</name>
        </tags>
        <tags id="1024">
            <name>1024</name>
        </tags>
        <tags id="1024">
            <name>1024</name>
        </tags>
    </group>
</root>`,
	},
}

func TestNew(t *testing.T) {
	a := assert.New(t)
	rslt := messagetest.NewMessageHandler()
	d := &ast.APIDoc{APIDoc: &ast.APIDocVersionAttribute{Value: xmlenc.String{Value: ast.Version}}}
	d.Parse(rslt.Handler, core.Block{Data: asttest.XML(a)})
	rslt.Handler.Stop()
	a.Empty(rslt.Errors)

	rslt = messagetest.NewMessageHandler()
	mock, err := New(rslt.Handler, d, indent, "/images", map[string]string{"client": "/test"}, testOptions)
	a.NotError(err).NotNil(mock)
	srv := rest.NewServer(t, mock, nil)

	// 测试路由是否正常
	srv.Post("/test/users", nil).Do().Status(http.StatusBadRequest)
	srv.Get("/test/users").Do().Status(http.StatusMethodNotAllowed)
	srv.Get("/not-found").Do().Status(http.StatusNotFound)
	srv.Get("/test/not-found").Do().Status(http.StatusNotFound)

	rslt.Handler.Stop()
	a.NotEmpty(rslt.Errors)
	srv.Close()

	rslt = messagetest.NewMessageHandler()
	mock, err = New(rslt.Handler, d, indent, "/images", map[string]string{"test": "/test"}, testOptions)
	a.NotError(err).NotNil(mock)
	srv = rest.NewServer(t, mock, nil)

	srv.Post("/test/users", nil).
		Header("accept", "application/json").
		Header("content-type", "application/xml").
		Body([]byte(`<root>
    <id>1</id>
    <name>n</name>
</root>`)).
		Do().
		Status(http.StatusCreated).
		Header("content-type", "application/json").
		BodyEmpty()

	// image
	imgBuffer := &bytes.Buffer{}
	srv.Get("/images/test.jpg").
		Header("accept", "image/Png").
		Do().
		Status(http.StatusOK).
		ReadBody(imgBuffer)
	img, err := png.Decode(imgBuffer)
	a.NotError(err).NotNil(img)
	a.Equal(img.Bounds(), image.Rect(0, 0, 500, 500))

	// image with size
	imgBuffer = &bytes.Buffer{}
	srv.Get("/images/test.jpg?width=200&height=50").
		Header("accept", "image/png;q=0.1,image/jpeg;q=0.9").
		Do().
		Status(http.StatusOK).
		ReadBody(imgBuffer)
	img, err = jpeg.Decode(imgBuffer)
	a.NotError(err).NotNil(img)
	a.Equal(img.Bounds(), image.Rect(0, 0, 200, 50))

	// image 不支持的 accept
	srv.Get("/images/test.jpg").
		Header("accept", "image/x-png").
		Do().
		Status(http.StatusNotAcceptable) // 无效的 accept

	// image 不在指定目录下，则 404
	srv.Get("/test.jpg").
		Header("accept", "image/png").
		Do().
		Status(http.StatusNotFound)

	rslt.Handler.Stop()
	a.Empty(rslt.Errors)

	// 版本号兼容性
	rslt = messagetest.NewMessageHandler()
	mock, err = New(rslt.Handler, &ast.APIDoc{APIDoc: &ast.APIDocVersionAttribute{Value: xmlenc.String{Value: "1.0.1"}}}, indent, "/images", nil, testOptions)
	a.Error(err).Nil(mock)
	rslt.Handler.Stop()
}

func TestLoad(t *testing.T) {
	a := assert.New(t)
	rslt := messagetest.NewMessageHandler()
	mock, err := Load(rslt.Handler, "./not-exists", indent, "/images", nil, testOptions)
	rslt.Handler.Stop()
	a.Error(err).Nil(mock)

	// LoadFromPath
	rslt = messagetest.NewMessageHandler()
	mock, err = Load(rslt.Handler, asttest.URI(a), indent, "/images", map[string]string{"admin": "/admin"}, testOptions)
	rslt.Handler.Stop()
	a.NotError(err).NotNil(mock)

	// loadFromURL
	static := http.FileServer(http.Dir(asttest.Dir(a)))
	srv := httptest.NewServer(static)
	defer srv.Close()

	rslt = messagetest.NewMessageHandler()
	mock, err = Load(rslt.Handler, core.URI(srv.URL+"/index.xml"), indent, "/images", map[string]string{"admin": "/admin"}, testOptions)
	rslt.Handler.Stop()
	a.NotError(err).NotNil(mock)
}

func TestIsValidRFC3339Date(t *testing.T) {
	a := assert.New(t)
	a.True(isValidRFC3339Date("2010-01-02"))
	a.False(isValidRFC3339Date("2010-01-32")) // 错误的日期
}

func TestIsValidRFC3339Time(t *testing.T) {
	a := assert.New(t)

	a.True(isValidRFC3339Time("17:18:19Z"))
	a.True(isValidRFC3339Time("17:18:19-08:30"))
	a.True(isValidRFC3339Time("17:18:19+08:30"))
	a.True(isValidRFC3339Time("17:18:19.22+08:30"))

	a.False(isValidRFC3339Time("17:18:19"))
	a.False(isValidRFC3339Time("17:18:69Z")) // 错误的日期
}

func TestIsValidRFC3339DateTime(t *testing.T) {
	a := assert.New(t)

	a.True(isValidRFC3339DateTime("2020-01-02T17:18:19Z"))
	a.True(isValidRFC3339DateTime("2020-01-02T17:18:19+08:00"))
	a.True(isValidRFC3339DateTime("2020-01-02T17:18:19-08:00"))

	a.False(isValidRFC3339DateTime("2020-01-02T17:18:19"))  // 错误的格式
	a.False(isValidRFC3339DateTime("2020-01-02T17:18:79Z")) // 错误的日期
	a.False(isValidRFC3339DateTime("2020-01-32T17:18:19Z")) // 错误的日期
}
