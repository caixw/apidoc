// SPDX-License-Identifier: MIT

package cmd

import (
	"flag"
	"testing"
	"time"

	"github.com/issue9/assert/v3"
)

var (
	_ flag.Getter = servers{}
	_ flag.Getter = &slice{}
	_ flag.Getter = &size{}
	_ flag.Getter = &dateRange{}
)

func TestServers_Set(t *testing.T) {
	a := assert.New(t, false)

	srv := servers{}
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

func TestSize_Set(t *testing.T) {
	a := assert.New(t, false)

	s := &size{}
	a.Error(s.Set(""))

	a.Error(s.Set(","))
	a.Error(s.Set("1,"))
	a.Error(s.Set(",5"))

	a.NotError(s.Set("1,5"))
	a.Equal(s.Min, 1).Equal(s.Max, 5)
}

func TestDateRange_Set(t *testing.T) {
	a := assert.New(t, false)

	d := &dateRange{}
	a.Error(d.Set(""))

	a.Error(d.Set(","))
	a.Error(d.Set("2020-01-02T17:04:15+01:00,"))
	a.Error(d.Set(",2020-01-07T17:04:15+01:00"))

	a.NotError(d.Set("2020-01-02T17:04:15+01:00,2020-01-05T17:04:15+01:00"))
	start, err := time.Parse(time.RFC3339, "2020-01-02T17:04:15+01:00")
	a.NotError(err)
	end, err := time.Parse(time.RFC3339, "2020-01-05T17:04:15+01:00")
	a.NotError(err)
	a.Equal(d.start, start).Equal(d.end, end)
}

func TestSlice_Set(t *testing.T) {
	a := assert.New(t, false)

	s := &slice{}
	a.NotError(s.Set("")).Equal(1, len(*s))
	a.NotError(s.Set(",")).Equal(2, len(*s))

}
