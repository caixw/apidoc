// SPDX-License-Identifier: MIT

package apidoc

import (
	"bytes"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v5/input"
	"github.com/caixw/apidoc/v5/internal/vars"
	"github.com/caixw/apidoc/v5/message"
)

func buildMessageHandle() (*bytes.Buffer, message.HandlerFunc) {
	buf := new(bytes.Buffer)

	return buf, func(msg *message.Message) {
		buf.WriteString(strconv.Itoa(int(msg.Type)))
		buf.WriteString(msg.Message)
	}
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

	out, f := buildMessageHandle()
	cfg := LoadConfig(message.NewHandler(f), wd)
	a.Empty(out.String()).NotNil(cfg)

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
	conf.Version = "5.0.1"
	err = conf.sanitize("./apidoc.yaml")
	a.Error(err).
		Equal(err.Field, "inputs")

	// 未声明 output
	conf.Inputs = []*input.Options{{}}
	err = conf.sanitize("./apidoc.yaml")
	a.Error(err).
		Equal(err.Field, "output")
}

func TestConfig_Do(t *testing.T) {
	a := assert.New(t)

	out, f := buildMessageHandle()
	h := message.NewHandler(f)
	LoadConfig(h, "./docs/example").Do(time.Now())
	a.Empty(out.String())
}

func TestConfig_Buffer(t *testing.T) {
	a := assert.New(t)

	out, f := buildMessageHandle()
	h := message.NewHandler(f)
	buf := LoadConfig(h, "./docs/example").Buffer()
	a.Empty(out.String()).
		True(buf.Len() > 0)
}
