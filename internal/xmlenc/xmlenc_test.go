// SPDX-License-Identifier: MIT

package xmlenc

import (
	"testing"

	"github.com/issue9/assert"
)

func TestBaseTag_SelfClose(t *testing.T) {
	a := assert.New(t)

	b := &BaseTag{}
	a.True(b.SelfClose())

	b = &BaseTag{EndTag: Name{
		Prefix: String{
			Value: "name",
		},
	}}
	a.True(b.SelfClose())

	b.EndTag.Local = String{
		Value: "name",
	}
	a.False(b.SelfClose())
}
