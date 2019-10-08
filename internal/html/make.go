// SPDX-License-Identifier: MIT

// +build ignore

package main

import (
	"bytes"
	"io/ioutil"
	"path/filepath"

	"github.com/issue9/utils"
)

const (
	header   = "// 由 make.go 生成，请勿修改！\n\n"
	pkgName  = "html"
	fileName = "static.go"
	resDir   = "../../docs" // 需要打包的资源文件所在的目录
)

type file struct {
	file        string // 相对于 resDir 的文件名
	contentType string
	name        string // 体现在 url 中的文件名部分
}

var files = []*file{
	{
		file:        "apidoc.xsl",
		contentType: "application/xslt+xml",
		name:        "apidoc.xsl",
	},
	{
		file:        "apidoc.css",
		contentType: "text/css",
		name:        "apidoc.css",
	},
	{
		file:        "apidoc.js",
		contentType: "application/javascript",
		name:        "apidoc.js",
	},
	{
		file:        "icon.svg",
		contentType: "image/svg+xml",
		name:        "icon.svg",
	},
}

func main() {
	buf := bytes.NewBufferString(header)
	buf.WriteString("package ")
	buf.WriteString(pkgName)
	buf.WriteByte('\n')

	defer func() {
		if err := utils.DumpGoFile(fileName, buf.String()); err != nil {
			panic(err)
		}
	}()

	buf.WriteString("var data=[]*static{\n")

	for _, file := range files {
		path := filepath.Join(resDir, file.file)
		if err := dump(buf, file, path); err != nil {
			panic(err)
		}
	}
	buf.WriteString("}")
}

func dump(buf *bytes.Buffer, file *file, path string) error {
	content, err := ioutil.ReadFile(path)

	ws := func(c string) {
		if err == nil {
			_, err = buf.WriteString(c)
		}
	}

	ws("{\n")

	ws("name:\"")
	ws(file.name)
	ws("\",\n")

	ws("contentType:\"")
	ws(file.contentType)
	ws("\",\n")

	ws("data:[]byte(`")
	ws(string(content))
	ws("`),\n")

	ws("},\n")
	return err
}
