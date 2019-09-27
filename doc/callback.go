// SPDX-License-Identifier: MIT

package doc

import (
	"encoding/xml"
	"strings"

	"github.com/caixw/apidoc/v5/internal/locale"
)

// HTTP 的安全模式
const (
	SchemaHTTP  = "HTTP"
	SchemaHTTPS = "HTTPS"
)

// Callback 回调函数的定义
//  <Callback deprecated="1.1.1" method="GET" schema="HTTPS">
//       <request status="200" mimetype="json" type="object">
//           <param name="name" type="string" />
//           <param name="sex" type="string">
//               <enum value="male">Male</enum>
//               <enum value="female">Female</enum>
//           </param>
//           <param name="age" type="number" />
//       </request>
//   </Callback>
type Callback struct {
	Schema      string     `xml:"schema,attr"` // http 或是 https，默认为 https
	Summary     string     `xml:"summary,attr,omitempty"`
	Description string     `xml:"description,omitempty"`
	Method      Method     `xml:"method,attr"`
	Queries     []*Param   `xml:"queries,omitempty"`
	Deprecated  Version    `xml:"deprecated,attr,omitempty"`
	Reference   string     `xml:"ref,attr,omitempty"`
	Responses   []*Request `xml:"response,omitempty"`
	Requests    []*Request `xml:"request"` // 至少一个
}

type shadowCallback Callback

// UnmarshalXML xml.Unmarshaler
func (c *Callback) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	name := "/" + start.Name.Local

	var cc shadowCallback
	if err := d.DecodeElement(&cc, &start); err != nil {
		return fixedSyntaxError(err, "", name, 0)
	}

	if cc.Method == "" {
		return newSyntaxError(name+"#method", locale.ErrRequired)
	}

	schema := strings.ToUpper(cc.Schema)
	if schema != SchemaHTTP && schema != SchemaHTTPS {
		return newSyntaxError(name+"#schema", locale.ErrInvalidValue)
	}

	if len(cc.Requests) == 0 {
		return newSyntaxError(name+"/request", locale.ErrRequired)
	}

	// 可以不需要 response

	*c = Callback(cc)
	return nil
}
