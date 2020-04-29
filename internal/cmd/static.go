// SPDX-License-Identifier: MIT

package cmd

import (
	"flag"
	"fmt"
	"io"
	"net/http"

	"github.com/caixw/apidoc/v7"
	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/locale"
)

var staticFlagSet *flag.FlagSet

var (
	staticPort        string
	staticDocs        uri
	staticStylesheet  bool
	staticContentType string
	staticURL         string
)

type uri core.URI

func (u uri) Get() interface{} {
	return string(u)
}

func (u *uri) Set(v string) error {
	fu, err := core.FileURI(v)
	if err != nil {
		return err
	}

	*u = uri(fu)
	return nil
}

func (u *uri) String() string {
	return core.URI(*u).String()
}

func initStatic() {
	staticFlagSet = command.New("static", static, staticUsage)
	staticFlagSet.StringVar(&staticPort, "p", ":8080", locale.Sprintf(locale.FlagStaticPortUsage))
	staticFlagSet.Var(&staticDocs, "docs", locale.Sprintf(locale.FlagStaticDocsUsage))
	staticFlagSet.StringVar(&staticContentType, "ct", "", locale.Sprintf(locale.FlagStaticContentTypeUsage))
	staticFlagSet.StringVar(&staticURL, "url", "", locale.Sprintf(locale.FlagStaticURLUsage))
	staticFlagSet.BoolVar(&staticStylesheet, "stylesheet", false, locale.Sprintf(locale.FlagStaticStylesheetUsage))
}

func static(io.Writer) (err error) {
	var path core.URI
	if staticFlagSet.NArg() != 0 {
		uri, err := getPath(staticFlagSet)
		if err != nil {
			return err
		}
		path = uri
	}

	h := core.NewMessageHandler(newHandlerFunc())
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
	h.Message(core.Succ, locale.ServerStart, url)

	return http.ListenAndServe(staticPort, handler)
}

func staticUsage(w io.Writer) error {
	_, err := fmt.Fprintln(w, locale.Sprintf(locale.CmdStaticUsage, getFlagSetUsage(staticFlagSet)))
	return err
}
