// SPDX-License-Identifier: MIT

package cmd

import (
	"io"
	"runtime"
	"strings"

	"github.com/caixw/apidoc/v5"
	"github.com/caixw/apidoc/v5/internal/locale"
	"github.com/caixw/apidoc/v5/internal/vars"
)

func initVersion() {
	command.New("version", version, buildUsage(locale.CmdVersionUsage))
}

func version(w io.Writer) error {
	goVersion := strings.TrimLeft(runtime.Version(), "go")
	printLocale(w, infoColor, locale.Version, apidoc.Version(), vars.DocVersion(), vars.CommitHash(), goVersion)
	return nil
}
