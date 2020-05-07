// SPDX-License-Identifier: MIT

// +build ignore

package main

import (
	"encoding/xml"
	"io/ioutil"
	"os"

	"github.com/caixw/apidoc/v7/internal/ast"
	"github.com/caixw/apidoc/v7/internal/docs"
	"github.com/caixw/apidoc/v7/internal/docs/makeutil"
	"github.com/caixw/apidoc/v7/internal/lang"
	"github.com/caixw/apidoc/v7/internal/vars"
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

func main() {
	for _, lang := range lang.Langs() {
		defaultConfig.Languages = append(defaultConfig.Languages, lang.DisplayName)
	}

	data, err := xml.MarshalIndent(defaultConfig, "", "\t")
	makeutil.PanicError(err)

	path, err := target.File()
	makeutil.PanicError(err)

	w := makeutil.NewWriter()
	w.WString(xml.Header).
		WString(fileHeader).
		WBytes(data).
		WString("\n") // 统一代码风格，文件末尾加一空行。

	makeutil.PanicError(ioutil.WriteFile(path, w.Bytes(), os.ModePerm))
}
