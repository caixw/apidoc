// SPDX-License-Identifier: MIT

// Package protocol 协议内容的定义
package protocol

// DocumentURI Many of the interfaces contain fields that correspond to the URI of a document.
// For clarity, the type of such a field is declared as a DocumentUri. Over the wire, it will still
// be transferred as a string, but this guarantees that the contents of that string can be parsed as a valid URI.
type DocumentURI string
