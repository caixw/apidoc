// SPDX-License-Identifier: MIT

package openapi

import (
	"testing"

	"github.com/issue9/assert"
)

func TestInfo_Sanitize(t *testing.T) {
	a := assert.New(t)

	info := &Info{}
	a.Error(info.Sanitize())

	info.Title = "title"
	a.Error(info.Sanitize())

	info.Title = "title"
	info.Version = "3.3.1"
	a.NotError(info.Sanitize())

	info.TermsOfService = "invalid url"
	a.Error(info.Sanitize())

	info.TermsOfService = "https://example.com"
	a.NotError(info.Sanitize())

	// contact

	info.Contact = &Contact{
		Name:  "name",
		Email: "invalid-email",
		URL:   "invalid-url",
	}
	a.Error(info.Sanitize())

	info.Contact.URL = "https://example.com"
	a.Error(info.Sanitize())

	info.Contact.Email = "user@example.com"
	a.NotError(info.Sanitize())

	// License
	info.License = &License{
		Name: "license",
		URL:  "invalid-url",
	}
	a.Error(info.Sanitize())

	info.License.URL = "https://example.com"
	a.NotError(info.Sanitize())
}
