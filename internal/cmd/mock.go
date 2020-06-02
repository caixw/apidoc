// SPDX-License-Identifier: MIT

package cmd

import (
	"bytes"
	"io"
	"net/http"
	"strings"

	"github.com/issue9/cmdopt"
	"github.com/issue9/rands"

	"github.com/caixw/apidoc/v7"
	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/locale"
)

// servers 参数
type servers map[string]string

func (s servers) Get() interface{} {
	return map[string]string(s)
}

func (s servers) Set(v string) error {
	pairs := strings.Split(v, ",")
	for _, pair := range pairs {
		index := strings.IndexByte(pair, '=')
		if index <= 0 {
			return locale.NewError(locale.ErrInvalidValue)
		}

		var v string
		if index < len(pair) {
			v = pair[index+1:]
		}
		s[strings.TrimSpace(pair[:index])] = v
	}

	return nil
}

func (s servers) String() string {
	if len(s) == 0 {
		return ""
	}

	buf := new(bytes.Buffer)

	for k, v := range s {
		buf.WriteString(k)
		buf.WriteByte('=')
		buf.WriteString(v)
		buf.WriteByte(',')
	}
	buf.Truncate(buf.Len() - 1)
	return buf.String()
}

var (
	mockPort        string
	mockServers     = make(servers, 0)
	mockStringAlpha string
	mockOptions     = &apidoc.MockOptions{}
	mockPath        = uri("./")
)

func initMock(command *cmdopt.CmdOpt) {
	mockFlagSet := command.New("mock", locale.Sprintf(locale.CmdMockUsage), doMock)
	mockFlagSet.StringVar(&mockPort, "p", ":8080", locale.Sprintf(locale.FlagMockPortUsage))
	mockFlagSet.Var(mockServers, "s", locale.Sprintf(locale.FlagMockServersUsage))
	mockFlagSet.Var(&mockPath, "path", locale.Sprintf(locale.FlagMockPathUsage))

	mockFlagSet.StringVar(&mockOptions.Indent, "indent", "\t", locale.Sprintf(locale.FlagMockIndentUsage))

	mockFlagSet.IntVar(&mockOptions.MaxSliceSize, "slice.max", 50, locale.Sprintf(locale.FlagMockSliceMaxUsage))
	mockFlagSet.IntVar(&mockOptions.MinSliceSize, "slice.min", 5, locale.Sprintf(locale.FlagMockSliceMinUsage))

	mockFlagSet.IntVar(&mockOptions.MaxNumber, "num.max", 10000, locale.Sprintf(locale.FlagMockNumMaxUsage))
	mockFlagSet.IntVar(&mockOptions.MinNumber, "num.min", 1, locale.Sprintf(locale.FlagMockNumMinUsage))
	mockFlagSet.BoolVar(&mockOptions.EnableFloat, "num.float", true, locale.Sprintf(locale.FlagMockNumFloatUsage))

	mockFlagSet.IntVar(&mockOptions.MaxString, "string.max", 64, locale.Sprintf(locale.FlagMockStringMaxUsage))
	mockFlagSet.IntVar(&mockOptions.MinString, "string.min", 24, locale.Sprintf(locale.FlagMockStringMinUsage))
	mockFlagSet.StringVar(&mockStringAlpha, "string.alpha", string(rands.AlphaNumber), locale.Sprintf(locale.FlagMockStringAlphaUsage))
}

func doMock(io.Writer) error {
	h := core.NewMessageHandler(messageHandle)
	defer h.Stop()

	mockOptions.Servers = mockServers
	mockOptions.StringAlpha = []byte(mockStringAlpha)
	handler, err := apidoc.MockFile(h, core.URI(mockPath), mockOptions)
	if err != nil {
		return err
	}

	url := "http://localhost" + mockPort
	h.Locale(core.Succ, locale.ServerStart, url)

	return http.ListenAndServe(mockPort, handler)
}
