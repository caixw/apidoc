// SPDX-License-Identifier: MIT

package lsp

import (
	"log"
	"os"

	"github.com/issue9/jsonrpc"
)

func newTestServer(header bool, info, erro *log.Logger) *server {
	return newServe(jsonrpc.NewStreamTransport(header, os.Stdin, os.Stdout, nil), info, erro)
}
