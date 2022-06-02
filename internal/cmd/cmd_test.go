// SPDX-License-Identifier: MIT

package cmd

import (
	"bytes"
	"flag"
	"testing"

	"github.com/issue9/assert/v2"
	"github.com/issue9/term/v3/colors"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/locale"
)

var (
	uuu             = uri("")
	_   flag.Getter = &uuu
)

func resetPrinters() (erro, warn, succ, info *bytes.Buffer) {
	erro = new(bytes.Buffer)
	warn = new(bytes.Buffer)
	succ = new(bytes.Buffer)
	info = new(bytes.Buffer)

	printers = map[core.MessageType]*printer{
		core.Erro: {
			out:    erro,
			color:  colors.Red,
			prefix: locale.ErrorPrefix,
		},
		core.Warn: {
			out:    warn,
			color:  colors.Cyan,
			prefix: locale.WarnPrefix,
		},
		core.Info: {
			out:    info,
			color:  colors.Default,
			prefix: locale.InfoPrefix,
		},
		core.Succ: {
			out:    succ,
			color:  colors.Green,
			prefix: locale.SuccessPrefix,
		},
	}

	return
}

func TestMessageHandle(t *testing.T) {
	a := assert.New(t, false)

	erro := new(bytes.Buffer)
	warn := new(bytes.Buffer)
	info := new(bytes.Buffer)
	succ := new(bytes.Buffer)

	printers[core.Erro].out = erro
	printers[core.Warn].out = warn
	printers[core.Info].out = info
	printers[core.Succ].out = succ

	h := core.NewMessageHandler(messageHandle)
	a.NotNil(h)

	h.Locale(core.Erro, "erro")
	h.Locale(core.Warn, "warn")
	h.Locale(core.Info, "info")
	h.Locale(core.Succ, "succ")
	h.Stop()

	a.Contains(erro.String(), "erro")
	a.Contains(warn.String(), "warn")
	a.Contains(info.String(), "info")
	a.Contains(succ.String(), "succ")
}
