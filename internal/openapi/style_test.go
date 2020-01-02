// SPDX-License-Identifier: MIT

package openapi

import (
	"testing"

	"github.com/issue9/assert"
)

func TestStyle_sanitize(t *testing.T) {
	a := assert.New(t)

	s := &Style{}
	a.Error(s.sanitize())

	s.Style = StyleDeepObject
	a.NotError(s.sanitize())

	s.Style = "invalid-value..."
	a.Error(s.sanitize())
}
