// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"bufio"
	"io/ioutil"
	"os"
)

// 指定产生的文件名。
const fileName = "static.go"

// 指定包名。
const packageName = "static"

// 文件头部的警告内容
const warning = "// 该文件由make.go自动生成，请勿手动修改！\n\n"

// 指定所有需要序列化的文件名。
var files = []string{
	"./style.css",
}

var templates = []string{
	"./main.html",
	"./doc.html",
}

func main() {
	f, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	w.WriteString(warning)

	// 输出包定义
	w.WriteString("package ")
	w.WriteString(packageName)
	w.WriteString("\n\n")

	// 输出files变量的整体。
	w.WriteString("var files=map[string][]byte{\n")
	for _, file := range files {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			panic(err)
		}

		w.WriteByte('"')
		w.WriteString(file)
		w.WriteString(`":`)
		w.WriteString("[]byte(`")
		w.Write(data)
		w.WriteString("`),")
	}
	w.WriteString("}\n")

	// 输出template变量的整体。
	w.WriteString("var Templates=map[string]string{\n")
	for _, file := range templates {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			panic(err)
		}

		w.WriteByte('"')
		w.WriteString(file)
		w.WriteString(`":`)
		w.WriteString("`")
		w.Write(data)
		w.WriteString("`,")
	}
	w.WriteString("}\n")

	if err = w.Flush(); err != nil {
		panic(err)
	}
}
