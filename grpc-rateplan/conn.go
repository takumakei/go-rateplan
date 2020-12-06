package rateplan

import (
	"io"
	"net"

	"github.com/takumakei/go-rateplan"
)

type Conn struct {
	net.Conn
	w io.Writer
}

func NewConn(conn net.Conn, plans rateplan.RatePlans) *Conn {
	return &Conn{
		Conn: conn,
		w:    rateplan.NewWriter(conn, plans),
	}
}

func (conn *Conn) Write(b []byte) (int, error) {
	return conn.w.Write(b)
}
