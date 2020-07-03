// SPDX-License-Identifier: MIT

package lsp

import (
	"golang.org/x/text/message"

	"github.com/caixw/apidoc/v7/internal/locale"
	"github.com/caixw/apidoc/v7/internal/lsp/protocol"
)

// window/showMessage
//
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#window_showMessage
func (s *server) windowShowMessage(t protocol.MessageType, messag string) error {
	return s.Notify("window/showMessage", &protocol.ShowMessageParams{
		Type:    t,
		Message: messag,
	})
}

func (s *server) windowShowInfoMessage(key message.Reference, v ...interface{}) error {
	return s.windowLogMessage(protocol.MessageTypeInfo, locale.Sprintf(key, v...))
}

func (s *server) windowShowLogMessage(key message.Reference, v ...interface{}) error {
	return s.windowLogMessage(protocol.MessageTypeLog, locale.Sprintf(key, v...))
}

func (s *server) windowShowWarnMessage(key message.Reference, v ...interface{}) error {
	return s.windowLogMessage(protocol.MessageTypeWarning, locale.Sprintf(key, v...))
}

func (s *server) windowShowErrorMessage(key message.Reference, v ...interface{}) error {
	return s.windowLogMessage(protocol.MessageTypeError, locale.Sprintf(key, v...))
}

// window/showMessageRequest
//
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#window_showMessageRequest
func (s *server) windowShowMessageRequest(t protocol.MessageType, actions []protocol.MessageActionItem, message string) (*protocol.MessageActionItem, error) {
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

func (s *server) windowShowInfoMessageRequest(actions []protocol.MessageActionItem, key message.Reference, v ...interface{}) (*protocol.MessageActionItem, error) {
	return s.windowShowMessageRequest(protocol.MessageTypeInfo, actions, locale.Sprintf(key, v...))
}

func (s *server) windowShowLogMessageRequest(actions []protocol.MessageActionItem, key message.Reference, v ...interface{}) (*protocol.MessageActionItem, error) {
	return s.windowShowMessageRequest(protocol.MessageTypeLog, actions, locale.Sprintf(key, v...))
}

func (s *server) windowShowWarnMessageRequest(actions []protocol.MessageActionItem, key message.Reference, v ...interface{}) (*protocol.MessageActionItem, error) {
	return s.windowShowMessageRequest(protocol.MessageTypeWarning, actions, locale.Sprintf(key, v...))
}

func (s *server) windowShowErrorMessageRequest(actions []protocol.MessageActionItem, key message.Reference, v ...interface{}) (*protocol.MessageActionItem, error) {
	return s.windowShowMessageRequest(protocol.MessageTypeError, actions, locale.Sprintf(key, v...))
}

// window/logMessage
//
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#window_logMessage
func (s *server) windowLogMessage(t protocol.MessageType, message string) error {
	return s.Notify("window/logMessage", &protocol.LogMessageParams{
		Type:    t,
		Message: message,
	})
}

func (s *server) windowLogInfoMessage(key message.Reference, v ...interface{}) error {
	return s.windowLogMessage(protocol.MessageTypeInfo, locale.Sprintf(key, v...))
}

func (s *server) windowLogLogMessage(key message.Reference, v ...interface{}) error {
	return s.windowLogMessage(protocol.MessageTypeLog, locale.Sprintf(key, v...))
}

func (s *server) windowLogWarnMessage(key message.Reference, v ...interface{}) error {
	return s.windowLogMessage(protocol.MessageTypeWarning, locale.Sprintf(key, v...))
}

func (s *server) windowLogErrorMessage(key message.Reference, v ...interface{}) error {
	return s.windowLogMessage(protocol.MessageTypeError, locale.Sprintf(key, v...))
}
