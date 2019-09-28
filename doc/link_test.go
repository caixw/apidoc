// SPDX-License-Identifier: MIT

package doc

import (
	"encoding/xml"
	"testing"

	"github.com/issue9/assert"
)

var (
	_ xml.Unmarshaler = &Link{}
	_ xml.Unmarshaler = &Contact{}
)

func TestLink_UnmarshalXML(t *testing.T) {
	a := assert.New(t)

	obj := &Link{
		Text: "text",
		URL:  "https://example.com",
	}
	str := `<Link url="https://example.com">text</Link>`

	data, err := xml.Marshal(obj)
	a.NotError(err).Equal(string(data), str)

	obj1 := &Link{}
	a.NotError(xml.Unmarshal([]byte(str), obj1))
	a.Equal(obj1, obj)

	obj1 = &Link{}
	a.NotError(xml.Unmarshal([]byte(str), obj1))
	str = `<Link url="url">text</Link>`
	a.Error(xml.Unmarshal([]byte(str), obj1))

	obj1 = &Link{}
	str = `<Link url="https://example.com"></Link>`
	a.Error(xml.Unmarshal([]byte(str), obj1))

	obj1 = &Link{}
	str = `<Link url="https://example.com"><Link>`
	a.Error(xml.Unmarshal([]byte(str), obj1))
}

func TestContact_UnmarshalXML(t *testing.T) {
	a := assert.New(t)

	obj := &Contact{
		Name:  "name",
		URL:   "https://example.com",
		Email: "name@example.com",
	}
	str := `<Contact name="name"><url>https://example.com</url><email>name@example.com</email></Contact>`

	data, err := xml.Marshal(obj)
	a.NotError(err).Equal(string(data), str)

	obj1 := &Contact{}
	a.NotError(xml.Unmarshal([]byte(str), obj1))
	a.Equal(obj1, obj)

	// 缺少 name
	obj1 = &Contact{}
	str = `<Contact name=""><url>https://example.com</url><email>name@example.com</email></Contact>`
	a.Error(xml.Unmarshal([]byte(str), obj1))

	// url 格式不正确
	obj1 = &Contact{}
	str = `<Contact name="name"><url>url</url><email>name@example.com</email></Contact>`
	a.Error(xml.Unmarshal([]byte(str), obj1))

	// email 格式不正确
	obj1 = &Contact{}
	str = `<Contact name="name"><url>https://example.com</url><email>email</email></Contact>`
	a.Error(xml.Unmarshal([]byte(str), obj1))

	// url 和 email 都不存在
	obj1 = &Contact{}
	str = `<Contact name="name"></Contact>`
	a.Error(xml.Unmarshal([]byte(str), obj1))

	// 少标签
	obj1 = &Contact{}
	str = `<Contact name="name">`
	a.Error(xml.Unmarshal([]byte(str), obj1))
}
