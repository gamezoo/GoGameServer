package ipc

import (
	. "connectorServer/tools"
	"connectorServer/tools/consul"
	myGprc "connectorServer/tools/grpc"
	"connectorServer/tools/stack"
	"google.golang.org/grpc"
	"io"
)

type ServerRecvHandle func(stream *Stream, msg *Req)
type StreamCloseHandle func()

type Stream struct {
	Ipc_TransferServer
	closeHandles []StreamCloseHandle
}

func (this *Stream) AddCloseHandle(closeHandle StreamCloseHandle) {
	this.closeHandles = append(this.closeHandles, closeHandle)
}

func (this *Stream) close() {
	for _, cb := range this.closeHandles {
		cb()
	}
	this.closeHandles = nil
}

type Server struct {
	serverRecvHandle ServerRecvHandle
}

func (this *Server) Transfer(stream Ipc_TransferServer) error {
	defer stack.PrintPanicStackError()

	s := &Stream{stream, []StreamCloseHandle{}}

	for {
		in, err := s.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			s.close()
			return err
		}
		go this.dealServerRecvHandle(s, in)
	}
}

func (this *Server) dealServerRecvHandle(stream *Stream, msg *Req) {
	defer stack.PrintPanicStackError()

	this.serverRecvHandle(stream, msg)
}

func InitServer(serverName string, serverRecvHandle ServerRecvHandle) (string, error) {
	serverPort, err := myGprc.InitServer(func(server *grpc.Server) {
		//注册处理模块
		RegisterIpcServer(server, &Server{
			serverRecvHandle: serverRecvHandle,
		})
	})

	//注册到服务
	err = consul.InitServer(serverName, serverPort)
	if err != nil {
		return "", err
	}
	INFO("join consul service...." + serverPort)

	return serverPort, nil
}