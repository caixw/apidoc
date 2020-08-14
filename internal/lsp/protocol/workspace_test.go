// SPDX-License-Identifier: MIT

package protocol

import (
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v7/core"
)

func TestWorkspaceFolder_Contains(t *testing.T) {
	a := assert.New(t)

	data := []*struct {
		folder, path string
		contained    bool
	}{
		{
			contained: true,
		},
		{
			folder:    "file:///root/p1",
			path:      "file:///root/p1/p2",
			contained: true,
		},
		{
			folder:    "file:///root/p1",
			path:      "/root/p1/p2",
			contained: true,
		},
		{
			folder:    "http://root/p1",
			path:      "http://root/p1/p2",
			contained: true,
		},
		{
			folder:    "https://root/p1",
			path:      "https://root/p1/p2",
			contained: true,
		},
		{
			folder: "http://root/p1",
			path:   "https://root/p1/p2",
		},
		{
			folder: "http:///root/p1",
			path:   "/root/p1/p2",
		},
		{
			folder: "file:///root/p1",
			path:   "file:///root",
		},
	}

	for _, item := range data {
		folder := WorkspaceFolder{URI: core.URI(item.folder)}
		a.Equal(folder.Contains(core.URI(item.path)), item.contained)
	}
}
