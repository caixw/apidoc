// SPDX-License-Identifier: MIT

package protocol

import (
	"testing"

	"github.com/issue9/assert"
)

func TestFileURI(t *testing.T) {
	a := assert.New(t)

	path := "/path/file"
	uri := FileURI(path)
	a.Equal(uri, fileScheme+"://"+path)
	file, err := uri.File()
	a.NotError(err).Equal(path, file)
}
