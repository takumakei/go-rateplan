package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	itecho "github.com/takumakei/go-rateplan/grpc-rateplan/internal/test/echo"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/examples/features/proto/echo"
)

func main() {
	app := &cli.App{
		Name:   "main",
		Action: action,
	}
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}

func action(c *cli.Context) error {
	lis, err := net.Listen("tcp", "127.0.0.1:12921")
	if err != nil {
		return err
	}
	defer lis.Close()

	server := grpc.NewServer()
	echo.RegisterEchoServer(server, itecho.EchoServer{})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)
		defer signal.Stop(sigint)
		for {
			select {
			case <-sigint:
				server.Stop()
			case <-ctx.Done():
				return
			}
		}
	}()

	return server.Serve(lis)
}
