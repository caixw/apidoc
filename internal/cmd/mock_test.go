// SPDX-License-Identifier: MIT

package cmd

import (
	"flag"
	"testing"

	"github.com/issue9/assert"
)

var _ flag.Getter = servers{}

func TestMockOptions(t *testing.T) {
	a := assert.New(t)

	srv := make(servers, 0)
	a.Equal(srv, srv.Get())
	a.Equal(0, len(srv))

	a.Error(srv.Set(""))
	a.Equal(0, len(srv))
	a.Equal(srv.String(), "")

	a.NotError(srv.Set("k1=v1"))
	a.Equal(1, len(srv))
	a.Equal(srv["k1"], "v1")
	a.NotEmpty(srv.String())

	a.NotError(srv.Set("k1=v1,k2=v2"))
	a.Equal(2, len(srv))
	a.Equal(srv["k1"], "v1")
	a.NotEmpty(srv.String())

	a.NotError(srv.Set("k1= v1, k2= v2"))
	a.Equal(2, len(srv))
	a.Equal(srv["k1"], " v1")
	a.Equal(srv["k2"], " v2")
	a.NotEmpty(srv.String())
}
