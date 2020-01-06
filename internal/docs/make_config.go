// SPDX-License-Identifier: MIT

// +build ignore

package main

import (
	"bufio"
	"encoding/xml"
	"os"
	"path/filepath"

	"github.com/caixw/apidoc/v5/internal/lang"
	"github.com/caixw/apidoc/v5/internal/vars"
)

const (
	fileHeader = `
<!--
	该文件由 /internal/docs/make.go 生成，请勿手动修改。

	主要包含了项目的一些基础配置项。
-->
`
)

var target = filepath.Join(vars.DocsDir(), "config.xml")

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
	Version:   vars.DocVersion(),
	Repo:      vars.RepoURL,
	URL:       vars.OfficialURL,
	Languages: make([]string, 0, len(lang.Langs())),
}

func main() {
	for _, lang := range lang.Langs() {
		defaultConfig.Languages = append(defaultConfig.Languages, lang.DisplayName)
	}

	data, err := xml.MarshalIndent(defaultConfig, "", "\t")
	if err != nil {
		panic(err)
	}

	file, err := os.Create(target)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			panic(err)
		}
	}()

	w := bufio.NewWriter(file)

	if _, err = w.WriteString(xml.Header); err != nil {
		panic(err)
	}

	if _, err = w.WriteString(fileHeader); err != nil {
		panic(err)
	}

	if _, err = w.Write(data); err != nil {
		panic(err)
	}

	// 统一代码风格，文件末尾加一空行。
	if _, err = w.WriteString("\n"); err != nil {
		panic(err)
	}

	if err = w.Flush(); err != nil {
		panic(err)
	}
}
