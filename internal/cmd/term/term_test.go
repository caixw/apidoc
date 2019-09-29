// SPDX-License-Identifier: MIT

package term

import (
	"bytes"
	"testing"

	"github.com/issue9/assert"
	"github.com/issue9/term/colors"

	"github.com/caixw/apidoc/v5/internal/locale"
	"github.com/caixw/apidoc/v5/message"
)

func TestNewHandlerFunc(t *testing.T) {
	a := assert.New(t)

	erro := new(bytes.Buffer)
	info := new(bytes.Buffer)
	def := colors.Default
	h := message.NewHandler(NewHandlerFunc(erro, erro, info, info, def, def, def, def))
	h.Message(message.Erro, locale.ErrRequired)
	h.Message(message.Warn, locale.ErrRequired)
	h.Message(message.Succ, locale.ErrRequired)
	h.Message(message.Info, locale.ErrRequired)
	h.Stop()
	a.NotEmpty(erro.String()).
		NotEmpty(info.String())
}

func TestLocale(t *testing.T) {
	a := assert.New(t)

	buf := new(bytes.Buffer)
	Locale(buf, colors.Default, locale.ErrRequired)
	a.Contains(buf.String(), locale.ErrRequired)
}
