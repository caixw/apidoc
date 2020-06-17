// SPDX-License-Identifier: MIT

package ast

import (
	"strconv"
	"testing"

	"github.com/issue9/assert"
	"github.com/issue9/version"
)

func TestVersion(t *testing.T) {
	a := assert.New(t)
	a.True(version.SemVerValid(Version))

	v := &version.SemVersion{}
	a.NotError(version.Parse(v, Version))
	major, err := strconv.Atoi(MajorVersion[1:])
	a.NotError(err)
	a.Equal(major, v.Major)
}

func TestTrimLeftSpace(t *testing.T) {
	a := assert.New(t)

	data := []*struct {
		input, output string
	}{
		{},
		{
			input:  `abc`,
			output: `abc`,
		},
		{
			input:  `  abc`,
			output: `abc`,
		},
		{
			input:  "  abc\n",
			output: "abc\n",
		},
		{ // 缩进一个空格
			input:  "  abc\n abc\n",
			output: " abc\nabc\n",
		},
		{ // 缩进一个空格
			input:  "\n  abc\n abc\n",
			output: "\n abc\nabc\n",
		},
		{ // 缩进格式不相同，不会有缩进
			input:  "\t  abc\n abc\n",
			output: "\t  abc\n abc\n",
		},

		{
			input:  "\t  abc\n\t abc\n\t xx\n",
			output: " abc\nabc\nxx\n",
		},
		{
			input:  "\t  abc\n\t abc\nxx\n",
			output: "\t  abc\n\t abc\nxx\n",
		},

		{ // 包含相同的 \t  内容
			input:  "\t  abc\n\t  abc\n\t  xx\n",
			output: "abc\nabc\nxx\n",
		},

		{ // 部分空格相同
			input:  "\t\t  abc\n\t  abc\n\t  xx\n",
			output: "\t  abc\n  abc\n  xx\n",
		},
	}

	for i, item := range data {
		output := trimLeftSpace(item.input)
		a.Equal(output, item.output, "not equal @ %d\nv1=%#v\nv2=%#v\n", i, output, item.output)
	}
}
