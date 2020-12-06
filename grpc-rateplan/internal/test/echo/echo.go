package echo

import (
	"context"
	"errors"
	"fmt"
	"io"

	"google.golang.org/grpc/examples/features/proto/echo"
)

type EchoServer struct {
	echo.UnimplementedEchoServer
}

// UnaryEcho is unary echo.
func (es EchoServer) UnaryEcho(_ context.Context, req *echo.EchoRequest) (*echo.EchoResponse, error) {
	return &echo.EchoResponse{Message: req.Message}, nil
}

// ServerStreamingEcho is server side streaming.
func (es EchoServer) ServerStreamingEcho(req *echo.EchoRequest, stream echo.Echo_ServerStreamingEchoServer) error {
	for i := 0; i < 1024; i++ {
		err := stream.Send(&echo.EchoResponse{Message: req.Message})
		if err != nil {
			return err
		}
	}
	return nil
}

// ClientStreamingEcho is client side streaming.
func (es EchoServer) ClientStreamingEcho(stream echo.Echo_ClientStreamingEchoServer) error {
	var m string
	for {
		req, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return err
		}
		m = req.Message
	}
	return stream.SendAndClose(&echo.EchoResponse{Message: m})
}

// BidirectionalStreamingEcho is bidi streaming.
func (es EchoServer) BidirectionalStreamingEcho(stream echo.Echo_BidirectionalStreamingEchoServer) error {
	var m string = "hello"
	go func() {
		fmt.Println("enter")
		defer fmt.Println("exit")
		for {
			req, err := stream.Recv()
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				return
			}
			m = req.Message
		}
	}()

	for i := 0; i < 1024; i++ {
		if err := stream.Send(&echo.EchoResponse{Message: m}); err != nil {
			return err
		}
	}

	return nil
}
