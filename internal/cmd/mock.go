// SPDX-License-Identifier: MIT

package cmd

import (
	"io"
	"net/http"
	"strings"

	"github.com/issue9/cmdopt"
	"github.com/issue9/errwrap"
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

	var buf errwrap.Buffer
	for k, v := range s {
		buf.WString(k).WByte('=').WString(v).WByte(',')
	}
	buf.Truncate(buf.Len() - 1)
	if buf.Err != nil {
		panic(buf.Err)
	}
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
	fs := command.New("mock", locale.Sprintf(locale.CmdMockUsage), doMock)
	fs.StringVar(&mockPort, "p", ":8080", locale.Sprintf(locale.FlagMockPortUsage))
	fs.Var(mockServers, "s", locale.Sprintf(locale.FlagMockServersUsage))
	fs.Var(&mockPath, "path", locale.Sprintf(locale.FlagMockPathUsage))

	fs.StringVar(&mockOptions.Indent, "indent", "\t", locale.Sprintf(locale.FlagMockIndentUsage))

	fs.IntVar(&mockOptions.MaxSliceSize, "slice.max", 10, locale.Sprintf(locale.FlagMockSliceMaxUsage))
	fs.IntVar(&mockOptions.MinSliceSize, "slice.min", 1, locale.Sprintf(locale.FlagMockSliceMinUsage))

	fs.IntVar(&mockOptions.MaxNumber, "num.max", 10000, locale.Sprintf(locale.FlagMockNumMaxUsage))
	fs.IntVar(&mockOptions.MinNumber, "num.min", 1, locale.Sprintf(locale.FlagMockNumMinUsage))
	fs.BoolVar(&mockOptions.EnableFloat, "num.float", false, locale.Sprintf(locale.FlagMockNumFloatUsage))

	fs.IntVar(&mockOptions.MaxString, "string.max", 64, locale.Sprintf(locale.FlagMockStringMaxUsage))
	fs.IntVar(&mockOptions.MinString, "string.min", 5, locale.Sprintf(locale.FlagMockStringMinUsage))
	fs.StringVar(&mockStringAlpha, "string.alpha", string(rands.AlphaNumber), locale.Sprintf(locale.FlagMockStringAlphaUsage))
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
