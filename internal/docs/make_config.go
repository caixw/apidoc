// SPDX-License-Identifier: MIT

// +build ignore

package main

import (
	"github.com/caixw/apidoc/v7/internal/ast"
	"github.com/caixw/apidoc/v7/internal/docs"
	"github.com/caixw/apidoc/v7/internal/docs/makeutil"
	"github.com/caixw/apidoc/v7/internal/lang"
	"github.com/caixw/apidoc/v7/internal/locale"
	"github.com/caixw/apidoc/v7/internal/token"
	"github.com/caixw/apidoc/v7/internal/vars"
)

var target = docs.Dir().Append("config.xml")

type config struct {
	XMLName struct{} `xml:"config"`

	Name      string     `xml:"name"`
	Version   string     `xml:"version"`
	Repo      string     `xml:"repo"`
	URL       string     `xml:"url"`
	Languages []language `xml:"languages>language"`
}

type language struct {
	ID   string `xml:"id,attr"`
	Name string `xml:",chardata"`
}

var defaultConfig = &config{
	Name:      vars.Name,
	Version:   ast.Version,
	Repo:      vars.RepoURL,
	URL:       vars.OfficialURL,
	Languages: make([]language, 0, len(lang.Langs())),
}

func main() {
	for _, lang := range lang.Langs() {
		l := language{ID: lang.ID, Name: lang.DisplayName}
		defaultConfig.Languages = append(defaultConfig.Languages, l)
	}
	makeutil.PanicError(makeutil.WriteXML(target, defaultConfig, "\t"))

	for _, tag := range locale.Tags() {
		types, err := token.NewTypes(&ast.APIDoc{}, tag)
		makeutil.PanicError(err)

		target := docs.Dir().Append("types." + tag.String() + ".xml")
		makeutil.PanicError(makeutil.WriteXML(target, types, "\t"))
	}
}
