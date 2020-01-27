// SPDX-License-Identifier: MIT

package jsonrpc

import (
	"context"
	"encoding/json"
	"reflect"
	"strconv"
	"sync"

	"github.com/issue9/autoinc"

	"github.com/caixw/apidoc/v6/internal/locale"
)

// Conn 连接对象，json-rpc 客户端和服务端是对等的，两者都使用 conn 初始化。
type Conn struct {
	stream  Streamer
	servers sync.Map
	autoinc *autoinc.AutoInc
}

type handler struct {
	f       reflect.Value
	in, out reflect.Type
}

// NewConn 声明新的 Conn 实例
func NewConn(stream Streamer) *Conn {
	return &Conn{
		stream:  stream,
		autoinc: autoinc.New(0, 1, 1000),
	}
}

// Register 注册一个新的服务
//
// f 为处理服务的函数，其原始为以下方式：
//  func(notify bool, params, result interface{}) error
// 其中 notify 表示是否为通知类型的请求；params 为用户请求的对象；
// result 为返回给用户的数据对象；error 则为处理出错是的返回值。
// params 和 result 必须为指针类型。
//
// 返回值表示是否添加成功，在已经存在相同值时，会添加失败。
func (conn *Conn) Register(method string, f interface{}) bool {
	if _, found := conn.servers.Load(method); found {
		return false
	}

	conn.servers.Store(method, newHandler(f))
	return true
}

var errType = reflect.TypeOf((*error)(nil)).Elem()

func newHandler(f interface{}) *handler {
	t := reflect.TypeOf(f)

	if t.Kind() != reflect.Func ||
		t.NumIn() != 3 ||
		t.In(0).Kind() != reflect.Bool ||
		t.In(1).Kind() != reflect.Ptr ||
		t.In(2).Kind() != reflect.Ptr ||
		!t.Out(0).Implements(errType) {
		panic("函数签名不正确")
	}

	in := t.In(1).Elem()
	if in.Kind() == reflect.Func || in.Kind() == reflect.Ptr || in.Kind() == reflect.Invalid {
		panic("函数签名不正确")
	}

	out := t.In(2).Elem()
	if out.Kind() == reflect.Func || out.Kind() == reflect.Ptr || out.Kind() == reflect.Invalid {
		panic("函数签名不正确")
	}

	return &handler{
		f:   reflect.ValueOf(f),
		in:  in,
		out: out,
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
		req.ID = strconv.FormatInt(conn.autoinc.MustID(), 10)
	}

	if err = conn.stream.Write(req); err != nil {
		return err
	}

	if notify {
		return nil
	}

	resp := &Response{}
	if err = conn.stream.Read(resp); err != nil {
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
	if err := conn.stream.Read(req); err != nil {
		return conn.writeError("", CodeParseError, err, nil)
	}

	f, found := conn.servers.Load(req.Method)
	if !found {
		return conn.writeError("", CodeMethodNotFound, locale.Errorf(locale.ErrInvalidValue), nil)
	}
	h := f.(*handler)

	notify := reflect.ValueOf(req.ID == "")
	in := reflect.New(h.in)
	if err := json.Unmarshal([]byte(*req.Params), in.Interface()); err != nil {
		return conn.writeError("", CodeParseError, err, nil)
	}

	out := reflect.New(h.out)
	if errVal := h.f.Call([]reflect.Value{notify, in, out}); !errVal[0].IsNil() {
		return conn.writeError("", CodeInternalError, errVal[0].Interface().(error), nil)
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
	return conn.stream.Write(resp)
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

	return conn.stream.Write(resp)
}
