// SPDX-License-Identifier: MIT

package cmd

import (
	"flag"
	"fmt"
	"io"
	"net/http"

	"github.com/caixw/apidoc/v5"
	"github.com/caixw/apidoc/v5/internal/locale"
	"github.com/caixw/apidoc/v5/message"
)

var staticFlagSet *flag.FlagSet

var staticPort string
var staticDocs string
var staticStylesheet bool

func initStatic() {
	staticFlagSet = command.New("static", static, staticUsage)
	staticFlagSet.StringVar(&staticPort, "p", ":8080", locale.Sprintf(locale.FlagStaticPortUsage))
	staticFlagSet.StringVar(&staticDocs, "docs", "", locale.Sprintf(locale.FlagStaticDocsUsage))
	staticFlagSet.BoolVar(&staticStylesheet, "stylesheet", false, locale.Sprintf(locale.FlagStaticStylesheetUsage))
}

func static(io.Writer) error {
	// path := getPath(staticFlagSet)

	h := message.NewHandler(newHandlerFunc())
	defer h.Stop()

	handler := apidoc.Static(staticDocs, staticStylesheet)
	return http.ListenAndServe(staticPort, handler)
}

func staticUsage(w io.Writer) error {
	_, err := fmt.Fprintln(w, locale.Sprintf(locale.CmdStaticUsage, getFlagSetUsage(staticFlagSet)))
	return err
}
