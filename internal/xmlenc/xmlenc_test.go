// SPDX-License-Identifier: MIT

package xmlenc

import (
	"testing"

	"github.com/issue9/assert/v3"

	"github.com/caixw/apidoc/v7/internal/locale"
)

func TestBaseTag_SelfClose(t *testing.T) {
	a := assert.New(t, false)

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

func TestBase_Usage(t *testing.T) {
	a := assert.New(t, false)

	b := Base{}
	a.Empty(b.Usage())

	b.UsageKey = locale.ErrInvalidUTF8Character
	a.Equal(b.Usage(), locale.Sprintf(locale.ErrInvalidUTF8Character))
}
