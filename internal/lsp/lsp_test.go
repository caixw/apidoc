// SPDX-License-Identifier: MIT

package lsp

import (
	"bytes"
	"context"
	"errors"
	"log"
	"net"
	"testing"
	"time"

	"github.com/issue9/assert/v3"
	"github.com/issue9/jsonrpc"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/lsp/protocol"
)

func TestServe(t *testing.T) {
	a := assert.New(t, false)
	a.Error(Serve(true, "not-exists-type", "", time.Second, nil, nil))
}

func TestServe_udp(t *testing.T) {
	a := assert.New(t, false)
	info := new(bytes.Buffer)
	erro := new(bytes.Buffer)
	infoLog := log.New(info, "[INFO]", 0)
	erroLog := log.New(erro, "[ERRO]", 0)
	srvExit := make(chan struct{}, 1)
	header := true

	go func() {
		err := Serve(header, "udp", ":8089", time.Second, infoLog, erroLog)
		a.True(errors.Is(err, context.Canceled))
		srvExit <- struct{}{}
	}()
	time.Sleep(500 * time.Millisecond) // 等待服务启动完成

	clientT, err := jsonrpc.NewUDPClientTransport(header, ":8089", "", time.Second)
	a.NotError(err).NotNil(clientT)
	client := jsonrpc.NewServer().NewConn(clientT, erroLog)
	clientCtx, clientCancel := context.WithCancel(context.Background())
	clientExit := make(chan struct{}, 1)
	go func() {
		client.Serve(clientCtx)
		clientExit <- struct{}{}
	}()
	time.Sleep(500 * time.Millisecond) // 等待服务启动完成

	initialize := make(chan struct{}, 1)
	client.Send("initialize", &protocol.InitializeParams{}, func(result *protocol.InitializeResult) error {
		a.Equal(result.ServerInfo.Name, core.Name)
		initialize <- struct{}{}
		return nil
	})

	a.Empty(erro.String())
	<-initialize
	initialized := make(chan struct{}, 1)
	client.Send("initialized", &protocol.InitializedParams{}, func(result *any) error {
		initialized <- struct{}{}
		return nil
	})

	<-initialized
	shutdown := make(chan struct{}, 1)
	client.Send("shutdown", nil, func(result *any) error {
		shutdown <- struct{}{}
		return nil
	})

	<-shutdown
	clientCancel()
	<-srvExit
	<-clientExit
}

func TestServe_tcp(t *testing.T) {
	a := assert.New(t, false)
	info := new(bytes.Buffer)
	erro := new(bytes.Buffer)
	infoLog := log.New(info, "[INFO]", 0)
	erroLog := log.New(erro, "[ERRO]", 0)
	srvExit := make(chan struct{}, 1)
	header := true

	go func() {
		err := Serve(header, "tcp", ":8089", time.Second, infoLog, erroLog)
		a.True(errors.Is(err, context.Canceled))
		srvExit <- struct{}{}
	}()
	time.Sleep(500 * time.Millisecond) // 等待服务启动完成

	udpConn, err := net.Dial("tcp", "127.0.0.1:8089")
	a.NotError(err)

	client := jsonrpc.NewServer().NewConn(jsonrpc.NewSocketTransport(header, udpConn, time.Second), erroLog)
	clientCtx, clientCancel := context.WithCancel(context.Background())
	clientExit := make(chan struct{}, 1)
	go func() {
		client.Serve(clientCtx)
		clientExit <- struct{}{}
	}()
	time.Sleep(500 * time.Millisecond) // 等待服务启动完成

	initialize := make(chan struct{}, 1)
	client.Send("initialize", &protocol.InitializeParams{}, func(result *protocol.InitializeResult) error {
		a.Equal(result.ServerInfo.Name, core.Name)
		initialize <- struct{}{}
		return nil
	})

	a.Empty(erro.String())
	<-initialize
	initialized := make(chan struct{}, 1)
	client.Send("initialized", &protocol.InitializedParams{}, func(result *any) error {
		initialized <- struct{}{}
		return nil
	})

	<-initialized
	shutdown := make(chan struct{}, 1)
	client.Send("shutdown", nil, func(result *any) error {
		shutdown <- struct{}{}
		return nil
	})

	<-shutdown
	clientCancel()
	<-srvExit
	<-clientExit
}
