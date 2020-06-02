// SPDX-License-Identifier: MIT

package cmd

import (
	"bytes"
	"flag"
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v7/core"
)

var (
	uuu             = uri("")
	_   flag.Getter = &uuu
)

func TestMessageHandle(t *testing.T) {
	a := assert.New(t)

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
