// SPDX-License-Identifier: MIT

package cmd

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/caixw/apidoc/v6"
	"github.com/caixw/apidoc/v6/internal/locale"
)

var lspFlagSet *flag.FlagSet

var lspPort string
var lspMode string

func initLSP() {
	lspFlagSet = command.New("lsp", doLSP, lspUsage)
	lspFlagSet.StringVar(&lspPort, "p", ":8080", locale.Sprintf(locale.FlagLSPPortUsage))
	lspFlagSet.StringVar(&lspMode, "m", "http", locale.Sprintf(locale.FlagLSPModeUsage))
}

func doLSP(io.Writer) error {
	conn := apidoc.LSP(log.New(erroOut, "", 0))

	switch strings.ToLower(lspMode) {
	case "http":
		return http.ListenAndServe(lspPort, conn)
	case "tcp":
	// TODO
	case "udp":
	// TODO
	case "websocket":
	// TODO
	default:
		return locale.Errorf(locale.ErrInvalidValue)
	}

	return nil
}

func lspUsage(w io.Writer) error {
	_, err := fmt.Fprintln(w, locale.Sprintf(locale.CmdLSPUsage, getFlagSetUsage(lspFlagSet)))
	return err
}
