// SPDX-License-Identifier: MIT

package xmlenc

import "github.com/caixw/apidoc/v7/core"

type (
	// Name 表示 XML 中的名称
	Name struct {
		core.Location
		Prefix String
		Local  String
	}

	// StartElement 表示 XML 的元素
	StartElement struct {
		core.Location
		Name       Name
		Attributes []*Attribute
		SelfClose  bool // 是否自闭合
	}

	// EndElement XML 的结束元素
	EndElement struct {
		core.Location
		Name Name
	}

	// Instruction 表示 XML 的指令
	Instruction struct {
		core.Location
		Name       String
		Attributes []*Attribute
	}

	// Attribute 表示 XML 属性
	Attribute struct {
		core.Location
		Name  Name
		Value String
	}

	// String 表示 XML 的字符串数据
	String struct {
		core.Location
		Value string
	}

	// CData 表示 XML 的 CDATA 数据
	CData struct {
		BaseTag
		Value String
	}

	// Comment 表示 XML 的注释
	Comment struct {
		core.Location
		Value String
	}
)

// Match 是否与 end 相匹配
func (s *StartElement) Match(end *EndElement) bool {
	return s.Name.Equal(end.Name)
}

// Equal 两个 name 是否相等
func (n Name) Equal(v Name) bool {
	return n.Prefix.Value == v.Prefix.Value &&
		n.Local.Value == v.Local.Value
}

// String fmt.Stringer
func (n Name) String() string {
	if n.Prefix.Value == "" {
		return n.Local.Value
	}
	return n.Prefix.Value + ":" + n.Local.Value
}
