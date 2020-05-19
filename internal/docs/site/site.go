// SPDX-License-Identifier: MIT

// Package site 用于生成网站内容
//
// 包括网站的基本信息，以及文档的翻译内容等。
package site

import (
	"golang.org/x/text/language/display"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/ast"
	"github.com/caixw/apidoc/v7/internal/docs/makeutil"
	"github.com/caixw/apidoc/v7/internal/lang"
	"github.com/caixw/apidoc/v7/internal/locale"
)

const (
	siteFilename = "site.xml" // 配置文件的文件名
	docBasename  = "locale."  // 翻译文档文件名的前缀部分，一般格式为 docBasename.{locale}.xml
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
	ID    string `xml:"id,attr"`
	Href  string `xml:"href,attr"`
	Title string `xml:"title,attr"`
	Doc   string `xml:"doc,attr"`
}

type doc struct {
	XMLName  struct{}   `xml:"locale"`
	Spec     []*spec    `xml:"spec>type"`
	Commands []*command `xml:"commands>command"`
	Config   []*item    `xml:"config>item"`
}

type spec struct {
	Name  string   `xml:"name,attr,omitempty"`
	Usage innerXML `xml:"usage,omitempty"`
	Items []*item  `xml:"item,omitempty"`
}

type innerXML struct {
	Text string `xml:",innerxml"`
}

type item struct {
	Name     string `xml:"name,attr"` // 变量名
	Type     string `xml:"type,attr"` // 变量的类型
	Array    bool   `xml:"array,attr"`
	Required bool   `xml:"required,attr"`
	Usage    string `xml:",innerxml"`
}

type command struct {
	Name  string `xml:"name,attr"`
	Usage string `xml:",innerxml"`
}

// Write 输出站点中所有需要自动生成的内容
func Write(target core.URI) error {
	site, d, err := gen()
	if err != nil {
		return err
	}

	if err := makeutil.WriteXML(target.Append(siteFilename), site, "\t"); err != nil {
		return err
	}

	for filename, dd := range d {
		if err := makeutil.WriteXML(target.Append(filename), dd, "\t"); err != nil {
			return err
		}
	}

	return nil
}

func gen() (*site, map[string]*doc, error) {
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
	docs := make(map[string]*doc, len(tags))

	for _, tag := range tags {
		locale.SetTag(tag)

		id := tag.String()
		docFilename := buildDocFilename(id)

		href := "index.xml"
		if id != locale.DefaultLocaleID {
			href = "index." + id + ".xml"
		}
		site.Locales = append(site.Locales, loc{
			ID:    id,
			Href:  href,
			Title: display.Self.Name(tag),
			Doc:   docFilename,
		})

		dd, err := genDoc()
		if err != nil {
			return nil, nil, err
		}
		docs[docFilename] = dd
	}

	return site, docs, nil
}

func genDoc() (*doc, error) {
	doc := &doc{}

	if err := doc.newCommands(); err != nil {
		return nil, err
	}
	if err := doc.newConfig(); err != nil {
		return nil, err
	}
	if err := doc.newSpec(&ast.APIDoc{}); err != nil {
		return nil, err
	}

	return doc, nil
}

func buildDocFilename(id string) string {
	return docBasename + id + ".xml"
}
