// SPDX-License-Identifier: MIT

package cmd

import (
	"fmt"
	"io"
	"runtime"
	"strings"

	"github.com/issue9/cmdopt"

	"github.com/caixw/apidoc/v7"
	"github.com/caixw/apidoc/v7/internal/locale"
	"github.com/caixw/apidoc/v7/internal/openapi"
)

var versionKind string

func initVersion(command *cmdopt.CmdOpt) {
	fs := command.New("version", locale.Sprintf(locale.CmdVersionUsage), version)
	fs.StringVar(&versionKind, "kind", "all", locale.FlagVersionKindUsage)
}

func version(w io.Writer) error {
	var msg string
	switch strings.ToLower(versionKind) {
	case "doc":
		msg = apidoc.DocVersion
	case "lsp":
		msg = apidoc.LSPVersion
	case "openapi":
		msg = openapi.LatestVersion
	case "apidoc":
		msg = apidoc.Version(true)
	default: // all 与 default 采取相同的输出
		goVersion := strings.TrimLeft(runtime.Version(), "go")
		msg = locale.Sprintf(locale.Version, apidoc.Version(true), apidoc.DocVersion, apidoc.LSPVersion, openapi.LatestVersion, goVersion)
	}
	_, err := fmt.Fprintln(w, msg)
	return err
}
