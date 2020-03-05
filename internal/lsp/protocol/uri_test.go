// SPDX-License-Identifier: MIT

package protocol

import (
	"testing"

	"github.com/issue9/assert"
)

func TestFileURI(t *testing.T) {
	a := assert.New(t)

	path := "/path/file"
	a.Equal(FileURI(path), fileScheme+path)

	path = "path/file"
	a.Equal(FileURI(path), fileScheme+"/"+path)
}
