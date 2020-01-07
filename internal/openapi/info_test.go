// SPDX-License-Identifier: MIT

package openapi

import (
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v6/doc"
)

func TestNewContact(t *testing.T) {
	a := assert.New(t)

	input := &doc.Contact{
		Email: "user@example.com",
		URL:   "https://example.com",
		Name:  "name",
	}

	output := newContact(input)
	a.Equal(output.Email, input.Email)

	output = newContact(nil)
	a.Nil(output)
}

func TestInfo_sanitize(t *testing.T) {
	a := assert.New(t)

	info := &Info{}
	a.Error(info.sanitize())

	info.Title = "title"
	a.Error(info.sanitize())

	info.Title = "title"
	info.Version = "3.3.1"
	a.NotError(info.sanitize())

	info.TermsOfService = "invalid url"
	a.Error(info.sanitize())

	info.TermsOfService = "https://example.com"
	a.NotError(info.sanitize())

	// contact

	info.Contact = &Contact{
		Name:  "name",
		Email: "invalid-email",
		URL:   "invalid-url",
	}
	a.Error(info.sanitize())

	info.Contact.URL = "https://example.com"
	a.Error(info.sanitize())

	info.Contact.Email = "user@example.com"
	a.NotError(info.sanitize())

	// License
	info.License = &License{
		Name: "license",
		URL:  "invalid-url",
	}
	a.Error(info.sanitize())

	info.License.URL = "https://example.com"
	a.NotError(info.sanitize())
}
