// SPDX-License-Identifier: MIT

package site

import (
	"github.com/caixw/apidoc/v7/internal/ast"
	"github.com/caixw/apidoc/v7/internal/docs"
	"github.com/caixw/apidoc/v7/internal/docs/makeutil"
)

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

func writeDoc(filename string) {
	doc := &doc{}
	makeutil.PanicError(doc.newCommands())
	makeutil.PanicError(doc.newConfig())
	makeutil.PanicError(doc.newSpec(&ast.APIDoc{}))
	makeutil.PanicError(makeutil.WriteXML(docs.Dir().Append(filename), doc, "\t"))
}
