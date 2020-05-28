// SPDX-License-Identifier: MIT

package build

import (
	"testing"
	"time"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/core/messagetest"
	"github.com/caixw/apidoc/v7/internal/docs"
)

func TestAllowConfigFilenames(t *testing.T) {
	a := assert.New(t)

	a.True(len(allowConfigFilenames) > 0)
	for _, name := range allowConfigFilenames {
		a.True(len(name) > 1)
	}
}

func TestLoadConfig(t *testing.T) {
	a := assert.New(t)

	rslt := messagetest.NewMessageHandler()
	cfg := LoadConfig(rslt.Handler, docs.Dir().Append("example"))
	rslt.Handler.Stop()
	a.NotNil(cfg).
		Empty(rslt.Errors).
		Empty(rslt.Successes)

	rslt = messagetest.NewMessageHandler()
	cfg = LoadConfig(rslt.Handler, docs.Dir()) // 不存在 apidoc 的配置文件
	rslt.Handler.Stop()
	a.Nil(cfg).
		NotEmpty(rslt.Errors).
		Empty(rslt.Successes)
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

func TestConfig_Save(t *testing.T) {
	a := assert.New(t)

	wd := core.FileURI("./")
	a.NotEmpty(wd)
	cfg, err := DetectConfig(wd, true)
	a.NotError(err).NotNil(cfg)
	a.NotError(cfg.Save(wd))
}

func TestConfig_Test(t *testing.T) {
	a := assert.New(t)

	rslt := messagetest.NewMessageHandler()
	cfg := LoadConfig(rslt.Handler, docs.Dir().Append("example"))
	a.NotNil(cfg)
	cfg.Test()

	rslt.Handler.Stop()
	a.Empty(rslt.Errors).
		NotEmpty(rslt.Successes) // 有成功提示
}

func TestConfig_Build(t *testing.T) {
	a := assert.New(t)

	rslt := messagetest.NewMessageHandler()
	cfg := LoadConfig(rslt.Handler, docs.Dir().Append("example"))
	a.NotNil(cfg)
	cfg.Build(time.Now())

	rslt.Handler.Stop()
	a.NotEmpty(rslt.Successes). // 有成功提示
					Empty(rslt.Errors)
}

func TestConfig_Buffer(t *testing.T) {
	a := assert.New(t)

	rslt := messagetest.NewMessageHandler()
	cfg := LoadConfig(rslt.Handler, docs.Dir().Append("example"))
	a.NotNil(cfg)

	buf := cfg.Buffer()
	rslt.Handler.Stop()
	a.Empty(rslt.Errors).
		Empty(rslt.Successes).
		True(buf.Len() > 0)
}
