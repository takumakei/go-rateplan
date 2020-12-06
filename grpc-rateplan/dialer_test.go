package rateplan_test

import (
	"context"
	"net"
	"strconv"
	"testing"
	"time"

	rateplan "github.com/takumakei/go-rateplan/grpc-rateplan"
	itecho "github.com/takumakei/go-rateplan/grpc-rateplan/internal/test/echo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/examples/features/proto/echo"
	"google.golang.org/grpc/status"
)

func TestDialer(t *testing.T) {
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	port := strconv.Itoa(lis.Addr().(*net.TCPAddr).Port)
	s := grpc.NewServer()
	defer s.Stop()
	echo.RegisterEchoServer(s, itecho.EchoServer{})
	go s.Serve(lis)

	dialer, err := rateplan.ParseDialer("32b@24h")
	if err != nil {
		t.Fatal(err)
	}

	conn, err := grpc.Dial(
		"127.0.0.1:"+port,
		grpc.WithInsecure(),
		grpc.WithContextDialer(dialer.DialContext),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	cli := echo.NewEchoClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	resp, err := cli.UnaryEcho(ctx, &echo.EchoRequest{Message: "abc"})
	if err != nil {
		if status.Code(err) != codes.DeadlineExceeded {
			t.Fatal(err)
		}
		t.Log("ok deadline exceeded")
		return
	}
	t.Log(resp)
}
