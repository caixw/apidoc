// SPDX-License-Identifier: MIT

package lsp

import (
	"github.com/caixw/apidoc/v6/internal/lsp/protocol"
)

// window/showMessage
//
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#window_showMessage
func (s *server) windowShowMessage(message string, t protocol.MessageType) error {
	return s.Send("window/showMessage", &protocol.ShowMessageParams{
		Type:    t,
		Message: message,
	}, nil)
}

// window/showMessageRequest
//
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#window_showMessageRequest
func (s *server) windowShowMessageRequest(message string, t protocol.MessageType, actions ...protocol.MessageActionItem) (*protocol.MessageActionItem, error) {
	out := &protocol.MessageActionItem{}
	in := &protocol.ShowMessageRequestParams{
		Type:    t,
		Message: message,
		Actions: actions,
	}

	err := s.Send("window/showMessageRequest", in, out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

// window/logMessage
//
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#window_logMessage
func (s *server) windowLogMessage(message string, t protocol.MessageType) error {
	return s.Send("window/logMessage", &protocol.ShowMessageParams{
		Type:    t,
		Message: message,
	}, nil)
}
