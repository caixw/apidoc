// SPDX-License-Identifier: MIT

package lsp

import (
	"fmt"
	"reflect"
	"sort"
	"unicode"

	"github.com/issue9/sliceutil"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/ast"
	"github.com/caixw/apidoc/v7/internal/lsp/protocol"
	"github.com/caixw/apidoc/v7/internal/node"
	"github.com/caixw/apidoc/v7/internal/xmlenc"
)

type tokenBuilder struct {
	tag, attr, value int // 可用的 token

	// 每个二级数组长度为 5，表示一组 semanticToken 数据。
	// 数据分别为 绝对行号，当前行的绝对起始位置，长度，以及 token 和 modifier。
	tokens [][]int
}

// textDocument/semanticTokens
func (s *server) textDocumentSemanticTokens(notify bool, in *protocol.SemanticTokensParams, out *protocol.SemanticTokens) error {
	f := s.findFolder(in.TextDocument.URI)
	if f == nil {
		return nil
	}

	f.parsedMux.RLock()
	defer f.parsedMux.RUnlock()

	out.Data = semanticTokens(f.doc, in.TextDocument.URI, 0, 1, 2)
	return nil
}

// tag 表示标签名的颜色值；
// attr 表示属性；
// value 表示属性的颜色值；
func semanticTokens(doc *ast.APIDoc, uri core.URI, tag, attr, value int) []int {
	b := &tokenBuilder{
		tag:    tag,
		attr:   attr,
		value:  value,
		tokens: make([][]int, 0, 100),
	}

	if doc.URI == uri {
		b.parse(reflect.ValueOf(doc), "APIs")
	}

	for _, api := range doc.APIs {
		matched := api.URI == uri || (api.URI == "" && doc.URI == uri)
		if !matched {
			continue
		}

		b.parse(reflect.ValueOf(api))
	}

	b.sort()
	return b.build()
}

// line 和 start 都为未计算的原始值
func (b *tokenBuilder) append(r core.Range, token int) {
	if r.End.Line == 0 && r.End.Character == 0 { // 未初始化的段被 node.RealValue 初始化成了零值，其长度必为 0
		return
	}

	l := r.End.Character - r.Start.Character
	if l < 0 { // 可能存在长度为 0 的，比如 default="" 值的长度为 0
		panic(fmt.Sprintf("无效的参数 range，其长度为 %d", l))
	}

	b.tokens = append(b.tokens, []int{r.Start.Line, r.Start.Character, l, token, 0})
}

func (b *tokenBuilder) build() []int {
	ret := make([]int, 0, 5*len(b.tokens))

	var line, start int
	for _, token := range b.tokens {
		currLine, currStart := token[0], token[1]
		if token[0] == line { // 同一行，start 取相对值
			token[1] -= start
		}
		token[0] -= line
		line, start = currLine, currStart

		ret = append(ret, token...)
	}

	return ret
}

// sort 排序内容，按从小到大
func (b *tokenBuilder) sort() {
	sort.SliceStable(b.tokens, func(i, j int) bool {
		ii := b.tokens[i]
		jj := b.tokens[j]
		return ii[0] < jj[0] || (ii[0] == jj[0] && ii[1] < jj[1])
	})
}

func (b *tokenBuilder) parse(v reflect.Value, exclude ...string) {
	v = node.RealValue(v)
	t := v.Type()
	if t.Kind() != reflect.Struct {
		return
	}

	b.parseAnonymous(v)

	for i := 0; i < t.NumField(); i++ {
		tf := t.Field(i)
		if tf.Anonymous ||
			unicode.IsLower(rune(tf.Name[0])) ||
			sliceutil.Count(exclude, func(i int) bool { return exclude[i] == tf.Name }) > 0 { // 需要过滤的字段
			continue
		}

		vf := node.RealValue(v.Field(i))
		if vf.Kind() == reflect.Array || vf.Kind() == reflect.Slice {
			for j := 0; j < vf.Len(); j++ {
				b.parse(vf.Index(j))
			}
		} else {
			b.parse(vf)
		}
	}
}

func (b *tokenBuilder) parseAnonymous(v reflect.Value) {
	t := v.Type()
	switch elem := v.Interface().(type) {
	case xmlenc.BaseTag:
		b.append(elem.StartTag.Range, b.tag)
		if !elem.SelfClose() {
			b.append(elem.EndTag.Range, b.tag)
		}
	case ast.CData:
		b.append(elem.StartTag.Range, b.tag)
		b.append(elem.EndTag.Range, b.tag)
	case ast.Content:
	case ast.Attribute:
		b.append(elem.AttributeName.Range, b.attr)
		b.append(elem.Value.Range, b.value)
	case ast.NumberAttribute:
		b.append(elem.AttributeName.Range, b.attr)
		b.append(elem.Value.Range, b.value)
	case ast.BoolAttribute:
		b.append(elem.AttributeName.Range, b.attr)
		b.append(elem.Value.Range, b.value)
	case ast.VersionAttribute:
		b.append(elem.AttributeName.Range, b.attr)
		b.append(elem.Value.Range, b.value)
	case ast.DateAttribute:
		b.append(elem.AttributeName.Range, b.attr)
		b.append(elem.Value.Range, b.value)
	case ast.MethodAttribute:
		b.append(elem.AttributeName.Range, b.attr)
		b.append(elem.Value.Range, b.value)
	case ast.StatusAttribute:
		b.append(elem.AttributeName.Range, b.attr)
		b.append(elem.Value.Range, b.value)
	case ast.TypeAttribute:
		b.append(elem.AttributeName.Range, b.attr)
		b.append(elem.Value.Range, b.value)
	case ast.APIDocVersionAttribute:
		b.append(elem.AttributeName.Range, b.attr)
		b.append(elem.Value.Range, b.value)
	default:
		for i := 0; i < t.NumField(); i++ {
			tf := t.Field(i)
			if !tf.Anonymous || unicode.IsLower(rune(tf.Name[0])) {
				continue
			}

			vf := node.RealValue(v.Field(i))
			if vf.Kind() == reflect.Array || vf.Kind() == reflect.Slice {
				for j := 0; j < vf.Len(); j++ {
					b.parseAnonymous(vf.Index(j))
				}
			} else {
				b.parseAnonymous(vf)
			}
		}
	}
}
