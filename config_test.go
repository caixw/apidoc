// SPDX-License-Identifier: MIT

package apidoc

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v6/input"
	"github.com/caixw/apidoc/v6/internal/docs"
	"github.com/caixw/apidoc/v6/internal/vars"
	"github.com/caixw/apidoc/v6/message/messagetest"
)

func TestLoadConfig(t *testing.T) {
	a := assert.New(t)

	erro, succ, h := messagetest.MessageHandler()
	cfg := LoadConfig(h, docs.Path("example"))
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

	cfg, err = loadFile("./", "./failed.yaml")
	a.Error(err).Nil(cfg)
}

func TestDetect_Load(t *testing.T) {
	a := assert.New(t)

	wd, err := filepath.Abs("./")
	a.NotError(err).NotEmpty(wd)
	a.NotError(Detect(wd, true))

	// 删除生成的文件
	defer func() {
		a.NotError(os.Remove("./.apidoc.yaml"))
	}()

	erro, _, h := messagetest.MessageHandler()
	cfg := LoadConfig(h, wd)
	a.Empty(erro.String()).NotNil(cfg)

	a.Equal(cfg.Version, vars.Version()).
		Equal(cfg.Inputs[0].Lang, "go")
}

func TestConfig_sanitize(t *testing.T) {
	a := assert.New(t)

	// 错误的版本号格式
	conf := &Config{}
	err := conf.sanitize("./apidoc.yaml")
	a.Error(err).
		Equal(err.Field, "version")

	// 与当前程序的版本号不兼容
	conf.Version = "1.0"
	err = conf.sanitize("./apidoc.yaml")
	a.Error(err).
		Equal(err.Field, "version")

	// 未声明 inputs
	conf.Version = "6.0.1"
	err = conf.sanitize("./apidoc.yaml")
	a.Error(err).
		Equal(err.Field, "inputs")

	// 未声明 output
	conf.Inputs = []*input.Options{{}}
	err = conf.sanitize("./apidoc.yaml")
	a.Error(err).
		Equal(err.Field, "output")
}

func TestConfig_Test(t *testing.T) {
	a := assert.New(t)

	erro, succ, h := messagetest.MessageHandler()
	cfg := LoadConfig(h, docs.Path("example"))
	a.NotNil(cfg)
	cfg.Test()

	h.Stop()
	a.Empty(erro.String()).
		NotEmpty(succ.String()) // 有成功提示
}

func TestConfig_Do(t *testing.T) {
	a := assert.New(t)

	erro, succ, h := messagetest.MessageHandler()
	cfg := LoadConfig(h, docs.Path("example"))
	a.NotNil(cfg)
	cfg.Build(time.Now())

	h.Stop()
	a.NotEmpty(succ.String()). // 有成功提示
					Empty(erro.String())
}

func TestConfig_Buffer(t *testing.T) {
	a := assert.New(t)

	erro, succ, h := messagetest.MessageHandler()
	cfg := LoadConfig(h, docs.Path("example"))
	a.NotNil(cfg)

	buf := cfg.Buffer()
	h.Stop()
	a.Empty(erro.String()).
		Empty(succ.String()).
		True(buf.Len() > 0)
}
