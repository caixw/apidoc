// SPDX-License-Identifier: MIT

package cmd

import (
	"io"
	"net/http"

	"github.com/issue9/cmdopt"

	"github.com/caixw/apidoc/v7"
	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/locale"
)

var (
	staticPort        string
	staticDocs        uri
	staticStylesheet  bool
	staticContentType string
	staticURL         string
	staticPath        uri
)

func initStatic(command *cmdopt.CmdOpt) {
	fs := command.New("static", locale.Sprintf(locale.CmdStaticUsage), static)
	fs.StringVar(&staticPort, "p", ":8080", locale.Sprintf(locale.FlagStaticPortUsage))
	fs.Var(&staticDocs, "docs", locale.Sprintf(locale.FlagStaticDocsUsage))
	fs.StringVar(&staticContentType, "ct", "", locale.Sprintf(locale.FlagStaticContentTypeUsage))
	fs.StringVar(&staticURL, "url", "", locale.Sprintf(locale.FlagStaticURLUsage))
	fs.BoolVar(&staticStylesheet, "stylesheet", false, locale.Sprintf(locale.FlagStaticStylesheetUsage))
	fs.Var(&staticPath, "path", locale.Sprintf(locale.FlagStaticPathUsage))
}

func static(io.Writer) (err error) {
	path := core.URI(staticPath)
	h := core.NewMessageHandler(messageHandle)
	defer h.Stop()

	var handler http.Handler

	if path == "" {
		handler = apidoc.Static(core.URI(staticDocs), staticStylesheet)
	} else {
		handler, err = apidoc.ViewFile(http.StatusOK, staticURL, path, staticContentType, core.URI(staticDocs), staticStylesheet)
		if err != nil {
			return err
		}
	}

	url := "http://localhost" + staticPort
	h.Locale(core.Succ, locale.ServerStart, url)

	return http.ListenAndServe(staticPort, handler)
}
