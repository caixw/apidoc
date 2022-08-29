// SPDX-License-Identifier: MIT

// Package ast 定义文档的抽象语法树
package ast

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"strings"
	"time"
	"unicode"

	"github.com/caixw/apidoc/v7/core"
)

const (
	// Version 文档规范的版本
	Version = "6.1.0"

	// MajorVersion 文档规范的主版本信息
	MajorVersion = "v6"

	// 文档允许的最小长度
	//
	// 文档都是以 <api 或是 <apidoc 开头的，所以最起码要等于此值。
	minSize = len("<api/>")
)

// 表示支持的各种数据类型
const (
	TypeNone     = "" // 空值表示不输出任何内容，仅用于 Request
	TypeBool     = "bool"
	TypeObject   = "object"
	TypeNumber   = "number"
	TypeString   = "string"
	TypeInt      = "number.int"
	TypeFloat    = "number.float"
	TypeEmail    = "string.email"
	TypeURL      = "string.url"
	TypeImage    = "string.image"     // 表示图片的 URL
	TypeDate     = "string.date"      // RFC3339 full-date
	TypeTime     = "string.time"      // RFC3339 full-time
	TypeDateTime = "string.date-time" // RFC3339 full-date + full-time
)

// 富文本可用的类型
const (
	RichtextTypeHTML     = "html"
	RichtextTypeMarkdown = "markdown"
)

// 几种与时间类型相关的格式
const (
	DateFormat     = "2006-01-02"     // 对应 TypeDate
	TimeFormat     = "15:04:05Z07:00" // 对应 TypeTime
	DateTimeFormat = time.RFC3339     // 对应 TypeDateTime
)

type (
	// Number 表示 XML 的数值类型
	Number struct {
		core.Range
		Int     int
		Float   float64
		IsFloat bool
	}

	// Bool 表示 XML 的布尔值类型
	Bool struct {
		core.Range
		Value bool
	}

	// Date 日期类型
	Date struct {
		core.Range
		Value time.Time
	}

	// Reference 指向引用父对象的数据
	Reference struct {
		core.Location
		Target any // 引用当前父对象的实际数据
	}

	// Definition 指向父对象的实际定义数据
	Definition struct {
		core.Location
		Target any // 父对象数据定义内容
	}

	// Referencer 包含了 Reference 数据需要实现的接口
	Referencer interface {
		core.Searcher
		References() []*Reference
	}

	// Definitioner 包含了 Definition 数据需要实现的接口
	Definitioner interface {
		core.Searcher
		Definition() *Definition
	}
)

// ParseType 获取类型字符串中的原始类型和子类型
func ParseType(t string) (primitive, sub string) {
	index := strings.IndexByte(t, '.')
	if index == -1 {
		return t, ""
	}
	return t[:index], t[index+1:]
}

func trimLeftSpace(v string) string {
	var min []byte // 找出的最小行首相同空格内容

	s := bufio.NewScanner(strings.NewReader(v))
	s.Split(bufio.ScanLines)
	for s.Scan() {
		line := s.Bytes()
		if len(bytes.TrimSpace(line)) == 0 { // 忽略空行
			continue
		}

		var index int
		for i, b := range line {
			if !unicode.IsSpace(rune(b)) {
				index = i
				break
			}
		}

		switch {
		case index == 0: // 当前行顶格
			return v
		case len(min) == 0: // 未初始化 min，且 index > 0
			min = make([]byte, index)
			copy(min, line[:index])
		default:
			min = getSamePrefix(min, line[:index])
		}
	}

	if len(min) == 0 {
		return v
	}

	buf := bufio.NewReader(strings.NewReader(v))
	ret := make([]byte, 0, buf.Size())
	for {
		line, err := buf.ReadBytes('\n')
		if bytes.HasPrefix(line, min) {
			line = line[len(min):]
		}
		ret = append(ret, line...)

		if errors.Is(err, io.EOF) {
			break
		}
	}

	return string(ret)
}

func getSamePrefix(v1, v2 []byte) []byte {
	l1, l2 := len(v1), len(v2)
	l := l1
	if l1 > l2 {
		l = l2
	}

	for i := 0; i < l; i++ {
		if v1[i] != v2[i] {
			return v1[:i]
		}
	}
	return v1[:l]
}
