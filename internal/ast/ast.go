// SPDX-License-Identifier: MIT

// Package ast 为 xml 服务的抽象语法树
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
)

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
