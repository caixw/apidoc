// SPDX-License-Identifier: MIT

package site

import (
	"testing"

	"github.com/issue9/assert/v2"

	"github.com/caixw/apidoc/v7/internal/lang"
	"github.com/caixw/apidoc/v7/internal/locale"
)

func TestGen(t *testing.T) {
	a := assert.New(t, false)

	site, docs, err := gen()
	a.NotError(err).
		NotNil(site).
		NotNil(docs)

	a.Equal(len(site.Languages), len(lang.Langs())).
		Equal(len(site.Locales), len(locale.Tags()))

	a.Equal(len(docs), len(locale.Tags()))

	defLocale := docs[buildDocFilename(locale.DefaultLocaleID)]
	for _, cmd := range defLocale.Commands {
		if cmd.Name == "build" {
			a.Contains(locale.Translate(locale.DefaultLocaleID, locale.CmdBuildUsage), cmd.Usage)
		}
	}

	for _, item := range defLocale.Config {
		a.Equal(item.Usage, locale.Translate(locale.DefaultLocaleID, "usage-config-"+item.Name))
	}
}
