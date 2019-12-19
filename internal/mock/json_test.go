// SPDX-License-Identifier: MIT

package mock

import (
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v5/doc"
)

func init() {
	test = true
}

type builder struct {
	input  *doc.Request
	output string
	err    bool
}

func testBuilder(a *assert.Assertion, builders []*builder) {
	for _, builder := range builders {
		data, err := buildJSON(builder.input)
		if builder.err {
			a.Error(err).Nil(data)
		} else {
			a.NotError(err).Equal(string(data), builder.output)
		}
	}
}

func TestBuildJSON(t *testing.T) {
	a := assert.New(t)

	data := []*builder{
		{
			input:  nil,
			output: "null",
		},
		{
			input:  &doc.Request{},
			output: "",
		},
		{ // array
			input: &doc.Request{
				Type:  doc.Bool,
				Array: true,
			},
			output: `[
    true,
    true,
    true,
    true,
    true
]`,
		},
		{
			input: &doc.Request{
				Type: doc.Bool,
			},
			output: "true",
		},
		{ // Object
			input: &doc.Request{
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
			output: `{
    "name": "1024",
    "id": 1024
}`,
		},

		{ // 各类型混合
			input: &doc.Request{
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
			output: `{
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

	testBuilder(a, data)
}
