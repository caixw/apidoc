// SPDX-License-Identifier: MIT

package cmd

import (
	"bytes"
	"flag"
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v6/internal/locale"
	"github.com/caixw/apidoc/v6/message"
)

func TestNewHandlerFunc(t *testing.T) {
	a := assert.New(t)

	erro := new(bytes.Buffer)
	warn := new(bytes.Buffer)
	info := new(bytes.Buffer)
	succ := new(bytes.Buffer)

	erroOut = erro
	warnOut = warn
	infoOut = info
	succOut = succ

	f := newHandlerFunc()
	a.NotNil(f)
	h := message.NewHandler(f)
	a.NotNil(h)

	h.Message(message.Erro, "erro")
	h.Message(message.Warn, "warn")
	h.Message(message.Info, "info")
	h.Message(message.Succ, "succ")
	h.Stop()

	a.Contains(erro.String(), "erro")
	a.Contains(warn.String(), "warn")
	a.Contains(info.String(), "info")
	a.Contains(succ.String(), "succ")
}

func TestBuildUsage(t *testing.T) {
	a := assert.New(t)
	w := new(bytes.Buffer)

	f := buildUsage(locale.CmdHelpUsage)
	a.NotNil(f)
	a.NotError(f(w))
	a.Equal(w.String(), locale.Sprintf(locale.CmdHelpUsage)+"\n")
}

func TestGetFlagSetUsage(t *testing.T) {
	a := assert.New(t)

	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	fs.Bool("b", false, "xx")
	output := fs.Output()
	a.NotEmpty(getFlagSetUsage(fs))
	a.Equal(output, fs.Output())
}
