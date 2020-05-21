// SPDX-License-Identifier: MIT

package messagetest

import (
	"testing"

	"github.com/issue9/assert"
)

func TestNewMessageHandler(t *testing.T) {
	a := assert.New(t)

	rslt := NewMessageHandler()
	a.NotNil(rslt).NotNil(rslt.Handler)
	rslt.Handler.Error("error")
	rslt.Handler.Stop()
	a.Equal(rslt.Errors[0], "error")

	rslt = NewMessageHandler()
	a.NotNil(rslt).NotNil(rslt.Handler)
	rslt.Handler.Info("info")
	rslt.Handler.Stop()
	a.Equal(rslt.Infos[0], "info")

	rslt = NewMessageHandler()
	a.NotNil(rslt).NotNil(rslt.Handler)
	rslt.Handler.Success("success")
	rslt.Handler.Stop()
	a.Equal(rslt.Successes[0], "success")

	rslt = NewMessageHandler()
	a.NotNil(rslt).NotNil(rslt.Handler)
	rslt.Handler.Warning("warn")
	rslt.Handler.Stop()
	a.Equal(rslt.Warns[0], "warn")
}
