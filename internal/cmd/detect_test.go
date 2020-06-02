// SPDX-License-Identifier: MIT

package cmd

import (
	"bytes"
	"testing"

	"github.com/issue9/assert"
	"gopkg.in/yaml.v2"

	"github.com/caixw/apidoc/v7/build"
	"github.com/caixw/apidoc/v7/core/messagetest"
	"github.com/caixw/apidoc/v7/internal/ast"
	"github.com/caixw/apidoc/v7/internal/docs"
)

func TestCmdDetect(t *testing.T) {
	a := assert.New(t)

	buf := new(bytes.Buffer)
	path := docs.Dir().Append("example")

	cmd := Init(buf)
	resetPrinters()
	err := cmd.Exec([]string{"detect", "-d", path.String()})
	a.NotError(err)
	cfg := &build.Config{}
	a.NotError(yaml.Unmarshal(buf.Bytes(), cfg))
	a.Equal(cfg.Version, ast.Version)

	cmd = Init(buf)
	resetPrinters()
	err = cmd.Exec([]string{"detect", "-d", path.String(), "-w"})
	a.NotError(err)
	rslt := messagetest.NewMessageHandler()
	defer rslt.Handler.Stop()
	cfg2 := &build.Config{}
	data := build.LoadConfig(rslt.Handler, path)
	a.Empty(rslt.Errors).NotNil(data)
	a.NotError(yaml.Unmarshal(buf.Bytes(), cfg2))
	a.Equal(cfg, cfg2)
}
