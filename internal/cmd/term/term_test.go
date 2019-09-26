// SPDX-License-Identifier: MIT

package term

import (
	"bytes"
	"strings"
	"testing"

	"github.com/issue9/assert"
	"github.com/issue9/term/colors"

	"github.com/caixw/apidoc/v5/internal/lang"
	"github.com/caixw/apidoc/v5/internal/locale"
	"github.com/caixw/apidoc/v5/message"
)

func TestLangs(t *testing.T) {
	a := assert.New(t)

	langs := Langs(3)
	a.Equal(len(langs), len(lang.Langs()))
	for _, l := range langs {
		cnt := strings.Count(l, strings.Repeat(" ", 3))
		a.True(cnt >= 2)
	}

	langs = Langs(10)
	a.Equal(len(langs), len(lang.Langs()))
	for _, l := range langs {
		cnt := strings.Count(l, strings.Repeat(" ", 10))
		a.True(cnt >= 2)
	}
}

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
