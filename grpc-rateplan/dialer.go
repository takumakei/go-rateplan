package rateplan

import (
	"context"
	"net"
	"net/url"
	"strings"

	"github.com/takumakei/go-rateplan"
)

type Dialer struct {
	Plans rateplan.RatePlans
}

func NewDialer(plans rateplan.RatePlans) *Dialer {
	return &Dialer{Plans: plans}
}

func ParseDialer(s string) (*Dialer, error) {
	plans, err := rateplan.ParseRatePlans(s)
	if err != nil {
		return nil, err
	}
	return NewDialer(plans), nil
}

func (d *Dialer) DialContext(ctx context.Context, target string) (net.Conn, error) {
	network, addr := parseDialTarget(target)
	conn, err := (&net.Dialer{}).DialContext(ctx, network, addr)
	if err != nil {
		return nil, err
	}
	return NewConn(conn, d.Plans), nil
}

// parseDialTarget returns the network and address to pass to dialer
func parseDialTarget(target string) (net string, addr string) {
	net = "tcp"

	m1 := strings.Index(target, ":")
	m2 := strings.Index(target, ":/")

	// handle unix:addr which will fail with url.Parse
	if m1 >= 0 && m2 < 0 {
		if n := target[0:m1]; n == "unix" {
			net = n
			addr = target[m1+1:]
			return net, addr
		}
	}
	if m2 >= 0 {
		t, err := url.Parse(target)
		if err != nil {
			return net, target
		}
		scheme := t.Scheme
		addr = t.Path
		if scheme == "unix" {
			net = scheme
			if addr == "" {
				addr = t.Host
			}
			return net, addr
		}
	}

	return net, target
}
