// SPDX-License-Identifier: MIT

package cmd

import (
	"fmt"
	"io"
	"runtime"
	"strings"

	"github.com/caixw/apidoc/v7"
	"github.com/caixw/apidoc/v7/internal/locale"
	"github.com/caixw/apidoc/v7/internal/openapi"
)

func initVersion() {
	command.New("version", locale.Sprintf(locale.CmdVersionUsage), version)
}

func version(w io.Writer) error {
	goVersion := strings.TrimLeft(runtime.Version(), "go")
	msg := locale.Sprintf(locale.Version, apidoc.Version(true), apidoc.DocVersion(), apidoc.LSPVersion(), openapi.LatestVersion, goVersion)
	_, err := fmt.Fprintln(w, msg)
	return err
}
