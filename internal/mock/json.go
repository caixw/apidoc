// SPDX-License-Identifier: MIT

package mock

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"

	"github.com/caixw/apidoc/v5/doc"
	"github.com/caixw/apidoc/v5/internal/locale"
	"github.com/caixw/apidoc/v5/message"
)

const indent = 4

type jsonBuilder struct {
	buf          *bytes.Buffer
	err          error
	deep         int
	indentString string
}

func validJSON(p *doc.Request, content []byte) error {
	if !json.Valid(content) {
		return message.NewLocaleError("", "request.body", 0, locale.ErrInvalidFormat)
	}

	r := bytes.NewReader(content)
	d := json.NewDecoder(r)

	for {
		token, err := d.Token()
		if err == io.EOF && token == nil { // 结束
			break
		}

		if err != nil {
			return err
		}

		if token == nil { // 对应 JSON null
			//
		}

		switch v := token.(type) {
		case string: // string
		case bool: // bool
		case json.Delim: // [、]、{、}
			switch v {
			case '[':
			case ']':
			case '{':
			case '}':
			}
		case float64: // number
		case json.Number: // number
		}
	}

	return nil
}

func (builder *jsonBuilder) writeIndent() *jsonBuilder {
	if builder.err != nil {
		return builder
	}
	_, builder.err = builder.buf.WriteString(builder.indentString)
	return builder
}

func (builder *jsonBuilder) incrIndent() *jsonBuilder {
	builder.deep++
	builder.indentString = strings.Repeat(" ", builder.deep*indent)
	return builder
}

func (builder *jsonBuilder) decrIndent() *jsonBuilder {
	builder.deep--
	builder.indentString = strings.Repeat(" ", builder.deep*indent)
	return builder
}

func (builder *jsonBuilder) writeStrings(str ...string) *jsonBuilder {
	if builder.err != nil {
		return builder
	}

	for _, s := range str {
		_, builder.err = builder.buf.WriteString(s)
		if builder.err != nil {
			break
		}
	}

	return builder
}

// v 只能是基本类型
func (builder *jsonBuilder) writeValue(v interface{}) *jsonBuilder {
	if builder.err != nil {
		return builder
	}

	vv, err := json.Marshal(v)
	if err != nil {
		builder.err = err
		return builder
	}

	_, builder.err = builder.buf.Write(vv)
	return builder
}

func buildJSON(p *doc.Request) ([]byte, error) {
	builder := &jsonBuilder{
		buf: new(bytes.Buffer),
	}

	if err := writeJSON(builder, p.ToParam(), true); err != nil {
		return nil, err
	}

	return builder.buf.Bytes(), nil
}

func writeJSON(builder *jsonBuilder, p *doc.Param, chkArray bool) error {
	if p == nil {
		builder.writeValue(nil)
		return builder.err
	}

	if p.Array && chkArray {
		builder.writeStrings("[\n").incrIndent()

		size := generateSliceSize()
		last := size - 1
		for i := 0; i < size; i++ {
			builder.writeIndent()

			if err := writeJSON(builder, p, false); err != nil {
				return err
			}

			if i < last {
				builder.writeStrings(",\n")
			} else {
				builder.writeStrings("\n")
			}
		}

		builder.decrIndent().writeIndent().writeStrings("]")
		return builder.err
	}

	switch p.Type {
	case doc.None:
		builder.writeValue(nil)
	case doc.Bool:
		builder.writeValue(generateBool())
	case doc.Number:
		builder.writeValue(generateNumber())
	case doc.String:
		builder.writeValue(generateString())
	case doc.Object:
		builder.writeStrings("{\n").incrIndent()

		last := len(p.Items) - 1
		for index, item := range p.Items {
			builder.writeIndent().writeStrings(`"`, item.Name, `"`, ": ")

			if err := writeJSON(builder, item, true); err != nil {
				return err
			}

			if index < last {
				builder.writeStrings(",\n")
			} else {
				builder.writeStrings("\n")
			}
		}

		builder.decrIndent().writeIndent().writeStrings("}")
	}

	return builder.err
}
