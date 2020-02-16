// SPDX-License-Identifier: MIT

package doc

import (
	"bytes"
	"testing"

	"github.com/issue9/assert"
)

func TestBegin(t *testing.T) {
	a := assert.New(t)

	a.True(bytes.HasPrefix([]byte("<apidoc "), apidocBegin))
	a.True(bytes.HasPrefix([]byte("<apidoc\t"), apidocBegin))
	a.True(bytes.HasPrefix([]byte("<apidoc\n"), apidocBegin))

	a.True(bytes.HasPrefix([]byte("<api "), apiBegin))
	a.True(bytes.HasPrefix([]byte("<api\t"), apiBegin))
	a.True(bytes.HasPrefix([]byte("<api\n"), apiBegin))
}
