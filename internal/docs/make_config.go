// SPDX-License-Identifier: MIT

// +build ignore

package main

import (
	"github.com/caixw/apidoc/v7/internal/ast"
	"github.com/caixw/apidoc/v7/internal/docs"
	"github.com/caixw/apidoc/v7/internal/docs/makeutil"
	"github.com/caixw/apidoc/v7/internal/lang"
	"github.com/caixw/apidoc/v7/internal/vars"
)

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

	makeutil.PanicError(makeutil.WriteXML(target, defaultConfig, "\t"))
}
