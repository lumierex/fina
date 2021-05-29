package frpc

import (
	"encoding/json"
	"fmt"
	"frpc/codec"
	"io"
	"log"
	"net"
	"reflect"
	"sync"
)

//
const MagicNumber = 0x3bef5c

// Option
// MagicNumber marks this is a rpc call
// CodeType client can choose different Code to encode body
type Option struct {
	MagicNumber int
	CodeType    codec.Type
}

var DefaultOption = &Option{
	MagicNumber: MagicNumber,
	CodeType:    codec.GobType,
}

// Custom protocol
// Options{MagicNumber: xxx, CodecType: xxx} | Header{ServiceMethod ...} | Body interface {}
// Once connection packets message maybe
// | Option | Header | Body | Option2 | Header2 | Body2

type Server struct {
}

func NewServer() *Server {
	return &Server{}
}

var DefaultServer = NewServer()

func Accept(lis net.Listener) {
	DefaultServer.Accept(lis)
}

func (server *Server) Accept(lis net.Listener) {
	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Println("rpc server: accept error: ", err)
			return
		}

		// TODO handle conn
		// read option
		// read header
		// read body
		go server.ServeConn(conn)

	}
}

func (server *Server) ServeConn(conn io.ReadWriteCloser) {
	// read Option
	// read header read body

	defer func() {
		_ = conn.Close()
	}()

	var opt Option
	if err := json.NewDecoder(conn).Decode(&opt); err != nil {
		log.Println("rpc server: read options error: ", err)
		return
	}

	if opt.MagicNumber != MagicNumber {
		log.Println("rpc server: invalid magic number: ", opt.MagicNumber)
		return
	}

	f := codec.NewCodecFuncMap[opt.CodeType]
	if f == nil {
		log.Println("rpc server: invalid code type: ", opt.CodeType)
		return
	}

	server.serveCodec(f(conn))
}

// placeholder for response argv when error occurs
var invalidRequest = struct{}{}

// serveCodec handle codec
// read request
// handle request
// send response
func (server *Server) serveCodec(codec codec.Codec) {
	// sending make sure send response complete
	sending := new(sync.Mutex)
	wg := new(sync.WaitGroup)

	for {
		req, err := server.readRequest(codec)
		if err != nil {
			// req is nil
			if req == nil {
				// can't recover connection, close connection
				break
			}
			req.h.Error = err.Error()
			// send request header with error back to client
			// handle request is async
			// response must be send on by one
			// so we use sending mutex to ensure
			server.sendResponse(codec, req.h, invalidRequest, sending)
			continue
		}
		wg.Add(1)
		go server.handleRequest(codec, req, sending, wg)
	}
	wg.Wait()
	_ = codec.Close()
}

// request stores all information of call
type request struct {
	h            *codec.Header
	argv, replyv reflect.Value
}

func (server *Server) readRequest(cc codec.Codec) (*request, error) {
	header, err := server.readRequestHeader(cc)
	if err != nil {
		return nil, err
	}
	req := &request{
		h: header,
	}

	// TODO read argv
	req.argv = reflect.New(reflect.TypeOf(""))
	if err = cc.ReadBody(req.argv.Interface()); err != nil {
		log.Println("rpc server: read argv err :", err)
	}
	return req, nil
}

func (server *Server) readRequestHeader(cc codec.Codec) (*codec.Header, error) {
	var h codec.Header
	if err := cc.ReadHeader(&h); err != nil {
		// TODO EOF write some
		if err != io.EOF && err != io.ErrUnexpectedEOF {
			log.Println("rpc server: read header error : ", err)
		}
		return nil, err
	}
	return &h, nil
}

func (server *Server) handleRequest(cc codec.Codec, req *request, sending *sync.Mutex, wg *sync.WaitGroup) {
	// go handler several request
	defer wg.Done()
	// record request header and body values
	log.Println(req.h, req.argv.Elem())
	req.replyv = reflect.ValueOf(fmt.Sprintf("frpc resp: %d", req.h.Seq))
	server.sendResponse(cc, req.h, req.replyv.Interface(), sending)
}

func (server *Server) sendResponse(cc codec.Codec, header *codec.Header, body interface{}, sending *sync.Mutex) {
	sending.Lock()
	defer sending.Unlock()

	if err := cc.Write(header, body); err != nil {
		log.Println("rpc server: write response error: ", err)
	}
}
