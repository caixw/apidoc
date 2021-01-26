// SPDX-License-Identifier: MIT

package cmd

import (
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/issue9/cmdopt"
	"github.com/issue9/errwrap"
	"github.com/issue9/rands"

	"github.com/caixw/apidoc/v7"
	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/locale"
)

// servers 参数
type (
	servers   map[string]string
	size      apidoc.Range
	slice     []string
	dateRange struct {
		start, end time.Time
	}
)

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

func (r size) Get() interface{} {
	return r
}

func (r *size) Set(v string) (err error) {
	pairs := strings.Split(v, ",")
	if len(pairs) != 2 {
		return locale.NewError(locale.ErrInvalidFormat)
	}

	left := strings.TrimSpace(pairs[0])
	if len(left) == 0 {
		return locale.NewError(locale.ErrInvalidFormat)
	}
	if r.Min, err = strconv.Atoi(left); err != nil {
		return err
	}

	right := strings.TrimSpace(pairs[1])
	if len(right) == 0 {
		return locale.NewError(locale.ErrInvalidFormat)
	}
	if r.Max, err = strconv.Atoi(right); err != nil {
		return err
	}

	return nil
}

func (r *size) String() string {
	return strconv.Itoa(r.Min) + "," + strconv.Itoa(r.Max)
}

func (s slice) Get() interface{} {
	return []string(s)
}

func (s *slice) Set(v string) (err error) {
	*s = strings.Split(v, ",")
	return nil
}

func (s *slice) String() string {
	return strings.Join(*s, ",")
}

func (d dateRange) Get() interface{} {
	return d
}

func (d *dateRange) Set(v string) (err error) {
	pairs := strings.Split(v, ",")
	if len(pairs) != 2 {
		return locale.NewError(locale.ErrInvalidFormat)
	}

	left := strings.TrimSpace(pairs[0])
	if len(left) == 0 {
		return locale.NewError(locale.ErrInvalidFormat)
	}
	if d.start, err = time.Parse(time.RFC3339, left); err != nil {
		return err
	}

	right := strings.TrimSpace(pairs[1])
	if len(right) == 0 {
		return locale.NewError(locale.ErrInvalidFormat)
	}
	if d.end, err = time.Parse(time.RFC3339, right); err != nil {
		return err
	}

	return nil
}

func (d *dateRange) String() string {
	return d.start.Format(time.RFC3339) + "," + d.end.Format(time.RFC3339)
}

var (
	mockOptions = &apidoc.MockOptions{}

	mockPort         string
	mockServers      = make(servers, 0)
	mockStringAlpha  string
	mockPath         = uri("./")
	mockSliceSize    = &size{Min: 5, Max: 10}
	mockNumberSize   = &size{Min: 100, Max: 10000}
	mockStringSize   = &size{Min: 50, Max: 1024}
	mockUsernameSize = &size{Min: 5, Max: 8}
	mockEmailDomains = &slice{"example.com"}
	mockURLDomains   = &slice{"https://example.com"}
	mockDateRange    = &dateRange{}
)

func initMock(command *cmdopt.CmdOpt) {
	fs := command.New("mock", locale.Sprintf(locale.CmdMockUsage), doMock)
	fs.StringVar(&mockPort, "p", ":8080", locale.Sprintf(locale.FlagMockPortUsage))
	fs.Var(mockServers, "servers", locale.Sprintf(locale.FlagMockServersUsage))
	fs.Var(&mockPath, "path", locale.Sprintf(locale.FlagMockPathUsage))

	fs.StringVar(&mockOptions.Indent, "indent", "\t", locale.Sprintf(locale.FlagMockIndentUsage))

	fs.Var(mockSliceSize, "slice.size", locale.Sprintf(locale.FlagMockSliceSizeUsage))

	fs.Var(mockNumberSize, "num.size", locale.Sprintf(locale.FlagMockNumSliceUsage))
	fs.BoolVar(&mockOptions.EnableFloat, "num.float", false, locale.Sprintf(locale.FlagMockNumFloatUsage))

	fs.Var(mockStringSize, "string.size", locale.Sprintf(locale.FlagMockStringSizeUsage))
	fs.StringVar(&mockStringAlpha, "string.alpha", string(rands.AlphaNumber), locale.Sprintf(locale.FlagMockStringAlphaUsage))

	fs.Var(mockUsernameSize, "email.username", locale.Sprintf(locale.FlagMockUsernameSizeUsage))
	fs.Var(mockEmailDomains, "email.domains", locale.Sprintf(locale.FlagMockEmailDomainsUsage))
	fs.Var(mockURLDomains, "url.domains", locale.Sprintf(locale.FlagMockURLDomainsUsage))

	fs.StringVar(&mockOptions.ImageBasePrefix, "image.prefix", "/__image__", locale.Sprintf(locale.FlagMockImagePrefixUsage))

	fs.Var(mockDateRange, "date.range", locale.Sprintf(locale.FlagMockDateRangeUsage))
}

func doMock(io.Writer) error {
	h := core.NewMessageHandler(messageHandle)
	defer h.Stop()

	mockOptions.Servers = mockServers
	mockOptions.StringAlpha = []byte(mockStringAlpha)
	mockOptions.SliceSize = apidoc.Range(*mockSliceSize)
	mockOptions.NumberSize = apidoc.Range(*mockNumberSize)
	mockOptions.StringSize = apidoc.Range(*mockStringSize)
	mockOptions.URLDomains = []string(*mockURLDomains)
	mockOptions.EmailDomains = []string(*mockEmailDomains)
	mockOptions.EmailUsernameSize = apidoc.Range(*mockUsernameSize)
	mockOptions.DateStart = mockDateRange.start
	mockOptions.DateEnd = mockDateRange.end
	handler, err := apidoc.MockFile(h, mockPath.URI(), mockOptions)
	if err != nil {
		return err
	}

	h.Locale(core.Succ, locale.ServerStart, mockPort)

	return http.ListenAndServe(mockPort, handler)
}
