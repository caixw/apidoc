// SPDX-License-Identifier: MIT

package doc

import (
	"testing"

	"github.com/issue9/assert"
)

func TestCheckXML(t *testing.T) {
	a := assert.New(t)

	xml := &XML{
		XMLAttr: true,
	}
	a.NotError(checkXML(xml, ""))

	xml.XMLExtract = true
	a.Error(checkXML(xml, ""))

	xml.XMLExtract = false
	xml.XMLNS = "https://example.com"
	a.Error(checkXML(xml, ""))

	xml.XMLExtract = false
	xml.XMLNS = ""
	xml.XMLNSPrefix = "ns"
	a.Error(checkXML(xml, ""))

	xml = &XML{
		XMLExtract: true,
	}
	a.NotError(checkXML(xml, ""))

	xml.XMLNS = "https://example.com"
	a.Error(checkXML(xml, ""))

	xml.XMLNS = ""
	xml.XMLNSPrefix = "ns"
	a.Error(checkXML(xml, ""))

	xml = &XML{}
	a.NotError(checkXML(xml, ""))

	xml.XMLNSPrefix = "ns"
	a.NotError(checkXML(xml, ""))

	xml.XMLNSPrefix = "ns"
	xml.XMLNS = "https://example.com"
	a.NotError(checkXML(xml, ""))

	xml.XMLNSPrefix = ""
	a.Error(checkXML(xml, ""))
}
