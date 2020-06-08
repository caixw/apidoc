// SPDX-License-Identifier: MIT

package cmd

import (
	"io"
	"log"

	"github.com/issue9/cmdopt"

	"github.com/caixw/apidoc/v7"
	"github.com/caixw/apidoc/v7/internal/locale"
)

var (
	lspPort   string
	lspMode   string
	lspHeader bool
)

func initLSP(command *cmdopt.CmdOpt) {
	ls := command.New("lsp", locale.Sprintf(locale.CmdLSPUsage), doLSP)
	ls.StringVar(&lspPort, "p", ":8080", locale.Sprintf(locale.FlagLSPPortUsage))
	ls.StringVar(&lspMode, "m", "http", locale.Sprintf(locale.FlagLSPModeUsage))
	ls.BoolVar(&lspHeader, "h", false, locale.Sprintf(locale.FlagLSPHeaderUsage))
}

func doLSP(o io.Writer) error {
	return apidoc.ServeLSP(lspHeader, lspMode, lspPort, log.New(o, "", 0), log.New(o, "", 0))
}
