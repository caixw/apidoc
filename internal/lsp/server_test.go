// SPDX-License-Identifier: MIT

package lsp

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/issue9/assert"
	"github.com/issue9/jsonrpc"

	"github.com/caixw/apidoc/v7/internal/lsp/protocol"
)

func newTestServer(header bool, info, erro *log.Logger) *server {
	return newServe(jsonrpc.NewStreamTransport(header, os.Stdin, os.Stdout, nil), info, erro)
}

func TestServer_setTrace(t *testing.T) {
	a := assert.New(t)
	s := newTestServer(true, log.New(ioutil.Discard, "", 0), log.New(ioutil.Discard, "", 0))

	err := s.setTrace(false, &protocol.SetTraceParams{}, nil)
	a.Error(err)
	jerr, ok := err.(*jsonrpc.Error)
	a.True(ok).Equal(jerr.Code, jsonrpc.CodeInvalidParams)
	a.Equal(s.trace, protocol.TraceValueOff)

	err = s.setTrace(false, &protocol.SetTraceParams{Value: protocol.TraceValueOff}, nil)
	a.NotError(err).Equal(s.trace, protocol.TraceValueOff)

	err = s.setTrace(false, &protocol.SetTraceParams{Value: protocol.TraceValueVerbose}, nil)
	a.NotError(err).Equal(s.trace, protocol.TraceValueVerbose)
}
