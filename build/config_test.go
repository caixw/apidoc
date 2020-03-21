// SPDX-License-Identifier: MIT

package build

import (
	"testing"
	"time"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v6/core"
	"github.com/caixw/apidoc/v6/core/messagetest"
	"github.com/caixw/apidoc/v6/internal/docs"
)

func TestLoadConfig(t *testing.T) {
	a := assert.New(t)

	erro, succ, h := messagetest.MessageHandler()
	cfg := LoadConfig(h, docs.Dir().Append("example"))
	h.Stop()
	a.NotNil(cfg).
		Empty(erro.String()).
		Empty(succ.String())

	erro, succ, h = messagetest.MessageHandler()
	cfg = LoadConfig(h, docs.Dir()) // 不存在 apidoc 的配置文件

	h.Stop()
	a.Nil(cfg).
		NotEmpty(erro.String()).
		Empty(succ.String())
}

func TestLoadFile(t *testing.T) {
	a := assert.New(t)

	cfg, err := loadFile("./", "./not-exists-file")
	a.Error(err).Nil(cfg)

	cfg, err = loadFile("./", "./testdata/failed.yaml")
	a.Error(err).Nil(cfg)
}

func TestConfig_sanitize(t *testing.T) {
	a := assert.New(t)

	// 错误的版本号格式
	conf := &Config{}
	err := conf.sanitize("./apidoc.yaml")
	err2, ok := err.(*core.SyntaxError)
	a.Error(err).
		True(ok).
		Equal(err2.Field, "version")

	// 与当前程序的版本号不兼容
	conf.Version = "1.0"
	err = conf.sanitize("./apidoc.yaml")
	err2, ok = err.(*core.SyntaxError)
	a.Error(err).
		True(ok).
		Equal(err2.Field, "version")

	// 未声明 inputs
	conf.Version = "6.0.1"
	err = conf.sanitize("./apidoc.yaml")
	err2, ok = err.(*core.SyntaxError)
	a.Error(err).
		True(ok).
		Equal(err2.Field, "inputs")

	// 未声明 output
	conf.Inputs = []*Input{{}}
	err = conf.sanitize("./apidoc.yaml")
	err2, ok = err.(*core.SyntaxError)
	a.Error(err).
		True(ok).
		Equal(err2.Field, "output")
}

func TestConfig_SaveToFile(t *testing.T) {
	a := assert.New(t)

	wd, err := core.FileURI("./")
	a.NotError(err).NotEmpty(wd)
	cfg, err := DetectConfig(wd, true)
	a.NotError(err).NotNil(cfg)
	a.NotError(cfg.SaveToFile(wd.Append(".apidoc.yaml")))
}

func TestConfig_Test(t *testing.T) {
	a := assert.New(t)

	erro, succ, h := messagetest.MessageHandler()
	cfg := LoadConfig(h, docs.Dir().Append("example"))
	a.NotNil(cfg)
	cfg.Test()

	h.Stop()
	a.Empty(erro.String()).
		NotEmpty(succ.String()) // 有成功提示
}

func TestConfig_Build(t *testing.T) {
	a := assert.New(t)

	erro, succ, h := messagetest.MessageHandler()
	cfg := LoadConfig(h, docs.Dir().Append("example"))
	a.NotNil(cfg)
	cfg.Build(time.Now())

	h.Stop()
	a.NotEmpty(succ.String()). // 有成功提示
					Empty(erro.String())
}

func TestConfig_Buffer(t *testing.T) {
	a := assert.New(t)

	erro, succ, h := messagetest.MessageHandler()
	cfg := LoadConfig(h, docs.Dir().Append("example"))
	a.NotNil(cfg)

	buf := cfg.Buffer()
	h.Stop()
	a.Empty(erro.String()).
		Empty(succ.String()).
		True(buf.Len() > 0)
}
