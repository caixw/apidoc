// SPDX-License-Identifier: MIT

package protocol

import (
	"testing"

	"github.com/issue9/assert"
)

func TestIsValidTraceValue(t *testing.T) {
	a := assert.New(t)
	a.True(IsValidTraceValue(TraceValueMessage))
	a.False(IsValidTraceValue("invalid-value"))
}

func TestBuildLogTrace(t *testing.T) {
	a := assert.New(t)

	p := BuildLogTrace(TraceValueOff, "m1", "v2")
	a.Nil(p)

	p = BuildLogTrace(TraceValueMessage, "m1", "v2")
	a.NotNil(p).Equal(p.Message, "m1").Empty(p.Verbose)

	p = BuildLogTrace(TraceValueVerbose, "m1", "v2")
	a.NotNil(p).Equal(p.Message, "m1").Equal(p.Verbose, "v2")

	a.Panic(func() {
		BuildLogTrace("invalid-trace", "m2", "v2")
	})
}
