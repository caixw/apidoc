// SPDX-License-Identifier: MIT

package doc

import (
	"encoding/xml"

	"github.com/caixw/apidoc/v5/internal/locale"
)

// Callback 回调函数的定义
//  <Callback deprecated="1.1.1" method="GET">
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
	Summary     string     `xml:"summary,attr,omitempty"`
	Path        *Path      `xml:"path,omitempty"`
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
	field := "/" + start.Name.Local

	shadow := (*shadowCallback)(c)
	if err := d.DecodeElement(shadow, &start); err != nil {
		return fixedSyntaxError(err, "", field, 0)
	}

	if shadow.Method == "" {
		return newSyntaxError(field+"#method", locale.ErrRequired)
	}

	if len(shadow.Requests) == 0 {
		return newSyntaxError(field+"/request", locale.ErrRequired)
	}

	// 可以不需要 response

	return nil
}
