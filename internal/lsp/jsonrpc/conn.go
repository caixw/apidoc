// SPDX-License-Identifier: MIT

package jsonrpc

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"reflect"
	"strconv"
	"sync"

	"github.com/caixw/apidoc/v6/internal/locale"
)

// ServerFunc 每个服务的函数签名
type ServerFunc func(in, out interface{}) error

// Conn 连接对象，json-rpc 客户端和服务端是对等的，两者都使用 conn 初始化。
type Conn struct {
	sequence int64
	errlog   *log.Logger
	stream   *stream
	servers  sync.Map
}

type handler struct {
	f       ServerFunc
	in, out reflect.Type
}

// NewConn 声明新的 Conn 实例
func NewConn(errlog *log.Logger, in io.Reader, out io.Writer) *Conn {
	return &Conn{
		errlog: errlog,
		stream: newStream(in, out),
	}
}

// Register 注册一个新的服务
//
// 返回值表示是否添加成功，在已经存在相同值时，会添加失败。
func (conn *Conn) Register(method string, f ServerFunc) bool {
	if _, found := conn.servers.Load(method); found {
		return false
	}

	conn.servers.Store(method, newHandler(f))
	return true
}

func newHandler(f ServerFunc) *handler {
	t := reflect.TypeOf(f)
	return &handler{
		f:   f,
		in:  t.Method(0).Type.Elem(),
		out: t.Method(0).Type.Elem(),
	}
}

// Notify 发送通知信息
func (conn *Conn) Notify(method string, in interface{}) error {
	return conn.send(true, method, in, nil)
}

// Send 发送请求内容，并获取其返回的数据
func (conn *Conn) Send(method string, in, out interface{}) error {
	return conn.send(false, method, in, out)
}

func (conn *Conn) send(notify bool, method string, in, out interface{}) error {
	data, err := json.Marshal(in)
	if err != nil {
		return err
	}

	req := &Request{
		Version: Version,
		Method:  method,
		Params:  (*json.RawMessage)(&data),
	}
	if !notify {
		req.ID = strconv.FormatInt(conn.sequence, 10)
	}

	if _, err = conn.stream.write(req); err != nil {
		return err
	}

	if notify {
		return nil
	}

	resp := &Response{}
	if err = conn.stream.readResponse(resp); err != nil {
		return err
	}

	if resp.Error != nil {
		return resp.Error
	}

	if req.ID != resp.ID {
		return NewError(CodeInvalidParams, locale.Sprintf(locale.VersionInCompatible))
	}

	return json.Unmarshal([]byte(*resp.Result), out)
}

// Serve 作为服务端运行
func (conn *Conn) Serve(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			if err := conn.serve(); err != nil {
				return err
			}
		}
	}
}

// 作为服务端，根据参数查找和执行服务
func (conn *Conn) serve() error {
	req := &Request{}
	if err := conn.stream.readRequest(req); err != nil {
		return conn.writeError("", CodeParseError, err, nil)
	}

	f, found := conn.servers.Load(req.Method)
	if !found {
		return conn.writeError("", CodeMethodNotFound, locale.Errorf(locale.ErrInvalidValue), nil)
	}
	h := f.(*handler)

	in := reflect.New(h.in).Interface()
	if err := json.Unmarshal([]byte(*req.Params), in); err != nil {
		return conn.writeError("", CodeParseError, err, nil)
	}

	out := reflect.New(h.out).Interface()
	if err := h.f(in, out); err != nil {
		return conn.writeError("", CodeInternalError, err, nil)
	}

	data, err := json.Marshal(out)
	if err != nil {
		return err
	}

	resp := &Response{
		Version: Version,
		Result:  (*json.RawMessage)(&data),
		ID:      req.ID,
	}
	_, err = conn.stream.write(resp)
	return err
}

func (conn *Conn) writeError(id string, code int, err error, data interface{}) error {
	resp := &Response{
		ID:      id,
		Version: Version,
	}

	if err2, ok := err.(*Error); ok {
		resp.Error = err2
	} else {
		resp.Error = NewError(code, err.Error())
	}

	_, err = conn.stream.write(resp)
	return err
}
