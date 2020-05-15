// SPDX-License-Identifier: MIT

// +build ignore

package main

import (
	"golang.org/x/text/language/display"

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
	Locales   []loc      `xml:"locales>locale"`
}

type language struct {
	ID   string `xml:"id,attr"`
	Name string `xml:",chardata"`
}

type loc struct {
	ID    string `xml:"id,attr"`
	Href  string `xml:"href,attr"`
	Title string `xml:"title,attr"`
	Types string `xml:"types,attr"`
}

var defaultConfig = &config{
	Name:      vars.Name,
	Version:   ast.Version,
	Repo:      vars.RepoURL,
	URL:       vars.OfficialURL,
	Languages: make([]language, 0, len(lang.Langs())),
	Locales:   make([]loc, 0, len(locale.Tags())),
}

func main() {
	for _, lang := range lang.Langs() {
		defaultConfig.Languages = append(defaultConfig.Languages, language{
			ID:   lang.ID,
			Name: lang.DisplayName,
		})
	}

	tags := locale.Tags()
	for _, tag := range tags {
		id := tag.String()
		href := "index.xml"
		if id != vars.DefaultLocaleID {
			href = "index." + id + ".xml"
		}
		defaultConfig.Locales = append(defaultConfig.Locales, loc{
			ID:    id,
			Href:  href,
			Title: display.Self.Name(tag),
			Types: "types." + id + ".xml",
		})
	}

	makeutil.PanicError(makeutil.WriteXML(target, defaultConfig, "\t"))

	for _, tag := range tags {
		types, err := token.NewTypes(&ast.APIDoc{}, tag)
		makeutil.PanicError(err)

		target := docs.Dir().Append("types." + tag.String() + ".xml")
		makeutil.PanicError(makeutil.WriteXML(target, types, "\t"))
	}
}
