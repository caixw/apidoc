// SPDX-License-Identifier: MIT

package cmd

import (
	"flag"
	"io"
	"runtime"
	"strings"

	"github.com/caixw/apidoc/v5"
	"github.com/caixw/apidoc/v5/internal/locale"
	"github.com/caixw/apidoc/v5/internal/vars"
)

var versionFlagSet *flag.FlagSet

func init() {
	versionFlagSet = command.New("version", version, versionUsage)
}

func version(w io.Writer) error {
	goVersion := strings.TrimLeft(runtime.Version(), "go")
	printLocale(w, infoColor, locale.Version, apidoc.Version(), vars.DocVersion(), vars.CommitHash(), goVersion)
	return nil
}

func versionUsage(w io.Writer) error {
	printLocale(w, infoColor, locale.FlagVUsage)
	return nil
}
