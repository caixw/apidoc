// SPDX-License-Identifier: MIT

package cmd

import (
	"flag"
	"fmt"
	"io"
	"net/http"

	"github.com/caixw/apidoc/v6"
	"github.com/caixw/apidoc/v6/internal/locale"
)

var lspFlagSet *flag.FlagSet

var lspPort string

func initLSP() {
	lspFlagSet = command.New("lsp", doLSP, lspUsage)
	lspFlagSet.StringVar(&lspPort, "p", ":8080", locale.Sprintf(locale.FlagLSPPortUsage))
}

func doLSP(io.Writer) error {
	handler, err := apidoc.LSP()
	if err != nil {
		return err
	}

	return http.ListenAndServe(lspPort, handler)
}

func lspUsage(w io.Writer) error {
	_, err := fmt.Fprintln(w, locale.Sprintf(locale.CmdLSPUsage, getFlagSetUsage(lspFlagSet)))
	return err
}
