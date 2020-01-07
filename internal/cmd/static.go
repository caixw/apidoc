// SPDX-License-Identifier: MIT

package cmd

import (
	"flag"
	"fmt"
	"io"
	"net/http"

	"github.com/caixw/apidoc/v6"
	"github.com/caixw/apidoc/v6/internal/locale"
	"github.com/caixw/apidoc/v6/message"
)

var staticFlagSet *flag.FlagSet

var (
	staticPort        string
	staticDocs        string
	staticStylesheet  bool
	staticContentType string
	staticURL         string
)

func initStatic() {
	staticFlagSet = command.New("static", static, staticUsage)
	staticFlagSet.StringVar(&staticPort, "p", ":8080", locale.Sprintf(locale.FlagStaticPortUsage))
	staticFlagSet.StringVar(&staticDocs, "docs", "", locale.Sprintf(locale.FlagStaticDocsUsage))
	staticFlagSet.StringVar(&staticContentType, "ct", "", locale.Sprintf(locale.FlagStaticContentTypeUsage))
	staticFlagSet.StringVar(&staticURL, "url", "", locale.Sprintf(locale.FlagStaticURLUsage))
	staticFlagSet.BoolVar(&staticStylesheet, "stylesheet", false, locale.Sprintf(locale.FlagStaticStylesheetUsage))
}

func static(io.Writer) error {
	path := getPath(staticFlagSet)

	h := message.NewHandler(newHandlerFunc())
	defer h.Stop()

	handler, err := apidoc.ViewFile(http.StatusOK, staticURL, path, staticContentType, staticDocs, staticStylesheet)
	if err != nil {
		return err
	}

	url := "http://localhost" + staticPort
	h.Message(message.Succ, locale.ServerStart, url)

	return http.ListenAndServe(staticPort, handler)
}

func staticUsage(w io.Writer) error {
	_, err := fmt.Fprintln(w, locale.Sprintf(locale.CmdStaticUsage, getFlagSetUsage(staticFlagSet)))
	return err
}
