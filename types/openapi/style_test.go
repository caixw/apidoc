// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package openapi

import (
	"testing"

	"github.com/issue9/assert"
)

func TestStyle_Sanitize(t *testing.T) {
	a := assert.New(t)

	s := &Style{}
	a.Error(s.Sanitize())

	s.Style = StyleDeepObject
	a.NotError(s.Sanitize())

	s.Style = "invalid-value..."
	a.Error(s.Sanitize())
}
