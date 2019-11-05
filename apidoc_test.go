// SPDX-License-Identifier: MIT

package apidoc

import (
	"bytes"
	"strconv"
	"testing"

	"github.com/issue9/assert"
	"github.com/issue9/version"

	"github.com/caixw/apidoc/v5/internal/vars"
	"github.com/caixw/apidoc/v5/message"
)

func buildMessageHandle() (*bytes.Buffer, message.HandlerFunc) {
	buf := new(bytes.Buffer)

	return buf, func(msg *message.Message) {
		buf.WriteString(strconv.Itoa(int(msg.Type)))
		buf.WriteString(msg.Message)
	}
}

func TestVersion(t *testing.T) {
	a := assert.New(t)

	a.Equal(Version(), vars.Version())
	a.True(version.SemVerValid(Version()))
}

func TestMake(t *testing.T) {
	a := assert.New(t)

	out, f := buildMessageHandle()
	h := message.NewHandler(f)
	Make(h, "./docs/example", true)
	a.Empty(out.Bytes())
}

func TestMakeBuffer(t *testing.T) {
	a := assert.New(t)

	out, f := buildMessageHandle()
	h := message.NewHandler(f)
	buf, dur := MakeBuffer(h, "./docs/example")
	a.Empty(out.Bytes()).
		True(dur > 0).
		True(buf.Len() > 0)
}
