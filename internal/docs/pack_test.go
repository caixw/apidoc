// SPDX-License-Identifier: MIT

package docs

import (
	"testing"

	"github.com/issue9/assert"
)

func TestGetFileInfos(t *testing.T) {
	a := assert.New(t)

	info, err := getFileInfos("./", nil)
	a.NotError(err).NotNil(info)
	a.Equal(8, len(info))
}
