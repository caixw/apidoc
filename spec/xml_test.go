// SPDX-License-Identifier: MIT

package spec

import (
	"testing"

	"github.com/issue9/assert"
)

func TestCheckXML(t *testing.T) {
	a := assert.New(t)

	xml := &XML{
		XMLAttr: true,
	}
	a.NotError(checkXML(false, false, xml, ""))

	xml.XMLExtract = true
	a.Error(checkXML(false, false, xml, ""))

	xml.XMLExtract = false
	xml.XMLNS = "https://example.com"
	a.Error(checkXML(false, false, xml, ""))

	xml.XMLExtract = false
	xml.XMLNS = ""
	xml.XMLNSPrefix = "ns"
	a.Error(checkXML(false, false, xml, ""))

	xml = &XML{
		XMLExtract: true,
	}
	a.NotError(checkXML(false, false, xml, ""))

	xml.XMLNS = "https://example.com"
	a.Error(checkXML(false, false, xml, ""))

	xml.XMLNS = ""
	xml.XMLNSPrefix = "ns"
	a.Error(checkXML(false, false, xml, ""))

	xml = &XML{}
	a.NotError(checkXML(false, false, xml, ""))

	xml.XMLNSPrefix = "ns"
	a.NotError(checkXML(false, false, xml, ""))

	xml.XMLNSPrefix = "ns"
	xml.XMLNS = "https://example.com"
	a.NotError(checkXML(false, false, xml, ""))

	xml.XMLNSPrefix = ""
	a.Error(checkXML(false, false, xml, ""))
}
