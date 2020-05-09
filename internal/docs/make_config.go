// SPDX-License-Identifier: MIT

// +build ignore

package main

import (
	l "golang.org/x/text/language"

	"github.com/caixw/apidoc/v7/internal/ast"
	"github.com/caixw/apidoc/v7/internal/docs"
	"github.com/caixw/apidoc/v7/internal/docs/makeutil"
	"github.com/caixw/apidoc/v7/internal/lang"
	"github.com/caixw/apidoc/v7/internal/token"
	"github.com/caixw/apidoc/v7/internal/vars"
)

var target = docs.Dir().Append("config.xml")
var typesTarget = docs.Dir().Append("types.cmn-hans.xml")

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
		l := language{ID: lang.Name, Name: lang.DisplayName}
		defaultConfig.Languages = append(defaultConfig.Languages, l)
	}
	makeutil.PanicError(makeutil.WriteXML(target, defaultConfig, "\t"))

	types, err := token.NewTypes(&ast.APIDoc{}, l.MustParse("cmn-hans"))
	makeutil.PanicError(err)
	makeutil.PanicError(makeutil.WriteXML(typesTarget, types, "\t"))
}
