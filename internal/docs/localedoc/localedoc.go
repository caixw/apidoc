// SPDX-License-Identifier: MIT

// Package localedoc 文档的本地化翻译内容
package localedoc

import "golang.org/x/text/language"

// LocaleDoc 管理本地化文档的集合
type LocaleDoc struct {
	XMLName  struct{}   `xml:"localedoc"`
	Types    []*Type    `xml:"types>type"`
	Commands []*Command `xml:"commands>command"`
}

// Type 用于生成文档中的类型信息
type Type struct {
	Name  string   `xml:"name,attr"`
	Usage InnerXML `xml:"usage"`
	Items []*Item  `xml:"item,omitempty"`
}

// InnerXML 可以用于在字符串嵌套 HTML
type InnerXML struct {
	Text string `xml:",innerxml"`
}

// Item 用于描述文档类型中的单条记录内容
type Item struct {
	Name     string `xml:"name,attr"` // 变量名
	Type     string `xml:"type,attr"` // 变量的类型
	Array    bool   `xml:"array,attr"`
	Required bool   `xml:"required,attr"`
	Usage    string `xml:",innerxml"`
}

// Command 命令行描述信息
type Command struct {
	Name  string `xml:"name,attr"`
	Usage string `xml:",innerxml"`
}

// Path 根据 tag 生成本地化的文件地址
func Path(tag language.Tag) string {
	return "localedoc." + tag.String() + ".xml"
}
