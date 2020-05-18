// SPDX-License-Identifier: MIT

// +build ignore

package main

import (
	"golang.org/x/text/language/display"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/ast"
	"github.com/caixw/apidoc/v7/internal/docs"
	"github.com/caixw/apidoc/v7/internal/docs/makeutil"
	"github.com/caixw/apidoc/v7/internal/lang"
	"github.com/caixw/apidoc/v7/internal/locale"
)

const target = "config.xml"

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
	ID       string `xml:"id,attr"`
	Href     string `xml:"href,attr"`
	Title    string `xml:"title,attr"`
	Types    string `xml:"types,attr"`
	Commands string `xml:"commands,attr"`
}

func main() {
	defaultConfig := &config{
		Name:      core.Name,
		Version:   ast.Version,
		Repo:      core.RepoURL,
		URL:       core.OfficialURL,
		Languages: make([]language, 0, len(lang.Langs())),
		Locales:   make([]loc, 0, len(locale.Tags())),
	}
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
		if id != locale.DefaultLocaleID {
			href = "index." + id + ".xml"
		}
		defaultConfig.Locales = append(defaultConfig.Locales, loc{
			ID:       id,
			Href:     href,
			Title:    display.Self.Name(tag),
			Types:    "types." + id + ".xml",
			Commands: "commands." + id + ".xml",
		})
	}

	makeutil.PanicError(makeutil.WriteXML(docs.Dir().Append(target), defaultConfig, "\t"))
}
