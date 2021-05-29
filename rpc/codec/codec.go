package codec

import (
	"io"
)

// err = client.Call("center.Hello", args, &reply)
// center => Service name
// Hello => method name

type Header struct {
	ServiceMethod string // center.Hello
	Seq           uint64 // sequence id to distinguish req from client
	Error         string
}

type Codec interface {
	io.Closer
	ReadHeader(*Header) error
	ReadBody(interface{}) error
	Write(*Header, interface{}) error
}

// return Codec constructor
// https://studygolang.com/articles/2271
type NewCodecFunc func(io.ReadWriteCloser) Codec
type Type string

const (
	GobType Type = "application/gob"
	//JsonType Type = "application/jso"
)

var NewCodecFuncMap map[Type]NewCodecFunc

func init() {
	NewCodecFuncMap = make(map[Type]NewCodecFunc)
	NewCodecFuncMap[GobType] = NewGobCodec
}
