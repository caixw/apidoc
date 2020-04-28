// SPDX-License-Identifier: MIT

// +build ignore

package main

import (
	"bufio"
	"encoding/xml"
	"os"

	"github.com/caixw/apidoc/v6/internal/ast"
	"github.com/caixw/apidoc/v6/internal/docs"
	"github.com/caixw/apidoc/v6/internal/lang"
	"github.com/caixw/apidoc/v6/internal/vars"
)

const fileHeader = "\n<!-- 该文件由工具自动生成，请勿手动修改！-->\n\n"

var target = docs.Dir().Append("config.xml")

type config struct {
	XMLName struct{} `xml:"config"`

	Name      string   `xml:"name"`
	Version   string   `xml:"version"`
	Repo      string   `xml:"repo"`
	URL       string   `xml:"url"`
	Languages []string `xml:"languages>language"`
}

var defaultConfig = &config{
	Name:      vars.Name,
	Version:   ast.Version,
	Repo:      vars.RepoURL,
	URL:       vars.OfficialURL,
	Languages: make([]string, 0, len(lang.Langs())),
}

func chkError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	for _, lang := range lang.Langs() {
		defaultConfig.Languages = append(defaultConfig.Languages, lang.DisplayName)
	}

	data, err := xml.MarshalIndent(defaultConfig, "", "\t")
	chkError(err)

	path, err := target.File()
	chkError(err)

	file, err := os.Create(path)
	chkError(err)
	defer func() {
		err = file.Close()
		chkError(err)
	}()

	w := bufio.NewWriter(file)

	_, err = w.WriteString(xml.Header)
	chkError(err)

	_, err = w.WriteString(fileHeader)
	chkError(err)

	_, err = w.Write(data)
	chkError(err)

	// 统一代码风格，文件末尾加一空行。
	_, err = w.WriteString("\n")
	chkError(err)

	err = w.Flush()
	chkError(err)
}
