// SPDX-License-Identifier: MIT

// +build ignore

package main

import (
	"bytes"
	"io/ioutil"

	"github.com/issue9/utils"
)

const (
	header  = "// 由 make.go 生成，请勿修改！\n\n"
	pkgName = "html"
	file    = "static.go"
)

var files = map[string]string{
	"XSL": "apidoc.xsl",
	"CSS": "apidoc.css",
	"JS":  "apidoc.js",
}

func main() {
	buf := bytes.NewBufferString(header)
	buf.WriteString("package ")
	buf.WriteString(pkgName)
	buf.WriteByte('\n')

	defer func() {
		if err := utils.DumpGoFile(file, buf.String()); err != nil {
			panic(err)
		}
	}()

	for name, file := range files {
		if err := dump(buf, name, file); err != nil {
			panic(err)
		}
	}
}

func dump(buf *bytes.Buffer, name, path string) error {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	buf.WriteString("var ")
	buf.WriteString(name)
	buf.WriteString("=[]byte(`")
	buf.Write(content)
	buf.WriteString("`)\n\n")
	return nil
}
