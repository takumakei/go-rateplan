package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	rateplan "github.com/takumakei/go-rateplan/grpc-rateplan"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/examples/features/proto/echo"
)

func main() {
	app := &cli.App{
		Name:   "main",
		Flags:  flags,
		Action: action,
	}
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}

var flags = []cli.Flag{
	flagRatePlan,
}

var flagRatePlan = &cli.StringFlag{
	Name:    "rateplan",
	Aliases: []string{"r"},
	Usage:   "`plan` e.g. 512kb@09:00/10s,512kb@09:01/10s",
}

func action(c *cli.Context) error {
	s := c.String(flagRatePlan.Name)
	if s == "" {
		now := time.Now()
		s = "8mb@" + now.Format("15:04:05-0700") + "/5s"
		s += ",8mb@" + now.Add(10*time.Second).Format("15:04:05-0700") + "/5s"
		s += ",8mb@" + now.Add(20*time.Second).Format("15:04:05-0700") + "/5s"
	}
	log.Println(s)

	dialer, err := rateplan.ParseDialer(s)
	if err != nil {
		return err
	}

	conn, err := grpc.Dial(
		"localhost:12921",
		grpc.WithInsecure(),
		grpc.WithContextDialer(dialer.DialContext),
	)
	if err != nil {
		return err
	}
	defer conn.Close()

	cli := echo.NewEchoClient(conn)

	ctx := context.Background()

	stream, err := cli.ClientStreamingEcho(ctx)
	if err != nil {
		return err
	}
	req := &echo.EchoRequest{Message: strings.Repeat("x", 1024)}
	a := time.Now()
	log.Println("begin", a.Format(time.RFC3339Nano))
	for i := 0; i < 1024; i++ {
		os.Stderr.Write([]byte{'.'})
		for j := 0; j < 4*1024; j++ {
			if err := stream.Send(req); err != nil {
				return err
			}
		}
	}
	b := time.Now()
	d := b.Sub(a)
	log.Println("begin", b.Format(time.RFC3339Nano), d)

	resp, err := stream.CloseAndRecv()
	if err != nil {
		return err
	}
	log.Println(resp.Message)
	return nil
}
