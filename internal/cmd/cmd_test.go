// SPDX-License-Identifier: MIT

package cmd

import (
	"bytes"
	"flag"
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/locale"
)

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

func TestGetPath(t *testing.T) {
	a := assert.New(t)
	fs := flag.NewFlagSet("test", flag.ContinueOnError)

	curr := core.FileURI("./")
	a.NotEmpty(curr)

	uri := getPath(fs)
	a.Equal(curr, uri)

	uri = getPath(nil)
	a.Equal(curr, uri)
}

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
