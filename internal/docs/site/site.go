// SPDX-License-Identifier: MIT

// Package localedoc 文档的本地化翻译内容
package site

import (
	"golang.org/x/text/language/display"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/ast"
	"github.com/caixw/apidoc/v7/internal/docs"
	"github.com/caixw/apidoc/v7/internal/docs/makeutil"
	"github.com/caixw/apidoc/v7/internal/lang"
	"github.com/caixw/apidoc/v7/internal/locale"
)

const (
	target      = "site.xml" // 配置文件的文件名
	docBasename = "locale."  // 翻译文档文件名的前缀部分，一般格式为 doc.{locale}.xml
)

type site struct {
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
	ID        string `xml:"id,attr"`
	Href      string `xml:"href,attr"`
	Title     string `xml:"title,attr"`
	LocaleDoc string `xml:"localedoc,attr"`
}

// Write 输出站点中所有需要自动生成的内容
func Write() {
	site := &site{
		Name:      core.Name,
		Version:   ast.Version,
		Repo:      core.RepoURL,
		URL:       core.OfficialURL,
		Languages: make([]language, 0, len(lang.Langs())),
		Locales:   make([]loc, 0, len(locale.Tags())),
	}
	for _, lang := range lang.Langs() {
		site.Languages = append(site.Languages, language{
			ID:   lang.ID,
			Name: lang.DisplayName,
		})
	}

	tags := locale.Tags()
	for _, tag := range tags {
		locale.SetTag(tag)

		id := tag.String()
		docFilename := docBasename + id + ".xml"

		href := "index.xml"
		if id != locale.DefaultLocaleID {
			href = "index." + id + ".xml"
		}
		site.Locales = append(site.Locales, loc{
			ID:        id,
			Href:      href,
			Title:     display.Self.Name(tag),
			LocaleDoc: docFilename,
		})

		writeDoc(docFilename)
	}

	makeutil.PanicError(makeutil.WriteXML(docs.Dir().Append(target), site, "\t"))
}
