// SPDX-License-Identifier: MIT

package build

import (
	"io/ioutil"
	"testing"

	"github.com/issue9/assert/v3"
	"gopkg.in/yaml.v3"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/core/messagetest"
	"github.com/caixw/apidoc/v7/internal/docs"
)

func TestAllowConfigFilenames(t *testing.T) {
	a := assert.New(t, false)

	a.True(len(allowConfigFilenames) > 0)
	for _, name := range allowConfigFilenames {
		a.True(len(name) > 1)
	}
}

func TestLoadConfig(t *testing.T) {
	a := assert.New(t, false)

	cfg, err := LoadConfig(docs.Dir().Append("example"))
	a.NotError(err).NotNil(cfg)

	// 读取的路径不应该是相对路径
	a.Equal(2, len(cfg.Inputs)).
		NotEqual(cfg.Inputs[0].Dir, ".").
		NotEqual(cfg.Inputs[1].Dir, ".").
		NotEqual(cfg.Output.Path, "./index.xml")

	cfg, err = LoadConfig(docs.Dir()) // 不存在 apidoc 的配置文件
	a.Error(err).Nil(cfg)

	// 不能加载远程配置项
	cfg, err = LoadConfig("https://apidoc.tools/example/.apidoc.yaml")
	a.Error(err).Nil(cfg)
}

func TestLoadFile(t *testing.T) {
	a := assert.New(t, false)

	cfg, err := loadFile("./", "./not-exists-file")
	a.Error(err).Nil(cfg)

	cfg, err = loadFile("./", "./testdata/failed.yaml")
	a.Error(err).Nil(cfg)
}

func TestConfig_sanitize(t *testing.T) {
	a := assert.New(t, false)

	// 错误的版本号格式
	conf := &Config{}
	err := conf.sanitize(".")
	err2, ok := err.(*core.Error)
	a.Error(err).
		True(ok).
		Equal(err2.Field, "version")

	// 与当前程序的版本号不兼容
	conf.Version = "1.0"
	err = conf.sanitize(".")
	err2, ok = err.(*core.Error)
	a.Error(err).
		True(ok).
		Equal(err2.Field, "version")

	// 未声明 inputs
	conf.Version = "6.0.1"
	err = conf.sanitize(".")
	err2, ok = err.(*core.Error)
	a.Error(err).
		True(ok).
		Equal(err2.Field, "inputs")

	// 未声明 output
	conf.Inputs = []*Input{{}}
	err = conf.sanitize(".")
	err2, ok = err.(*core.Error)
	a.Error(err).
		True(ok).
		Equal(err2.Field, "output")
}

func TestConfig_Save(t *testing.T) {
	a := assert.New(t, false)

	wd := core.FileURI("./")
	a.NotEmpty(wd)
	cfg, err := DetectConfig(wd, true)
	a.NotError(err).NotNil(cfg)
	a.NotError(cfg.Save(wd))

	// 通过 save 保存的路径应该是相对路径
	cfg = &Config{}
	data, err := ioutil.ReadFile("./" + allowConfigFilenames[0])
	a.NotError(err).NotNil(data)
	a.NotError(yaml.Unmarshal(data, cfg))
	a.Equal(".", cfg.Inputs[0].Dir)
	a.Equal("apidoc.xml", cfg.Output.Path)
}

func TestConfig_CheckSyntax(t *testing.T) {
	a := assert.New(t, false)

	cfg, err := LoadConfig(docs.Dir().Append("example"))
	a.NotError(err).NotNil(cfg)

	rslt := messagetest.NewMessageHandler()
	cfg.CheckSyntax(rslt.Handler)
	rslt.Handler.Stop()
	a.Empty(rslt.Errors)
}

func TestConfig_Build(t *testing.T) {
	a := assert.New(t, false)

	cfg, err := LoadConfig(docs.Dir().Append("example"))
	a.NotError(err).NotNil(cfg)

	rslt := messagetest.NewMessageHandler()
	cfg.Build(rslt.Handler)
	rslt.Handler.Stop()
	a.Empty(rslt.Errors)
}

func TestConfig_Buffer(t *testing.T) {
	a := assert.New(t, false)

	cfg, err := LoadConfig(docs.Dir().Append("example"))
	a.NotError(err).NotNil(cfg)

	rslt := messagetest.NewMessageHandler()
	buf := cfg.Buffer(rslt.Handler)
	rslt.Handler.Stop()
	a.Empty(rslt.Errors).
		Empty(rslt.Successes).
		True(buf.Len() > 0)
}
