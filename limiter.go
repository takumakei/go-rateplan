package rateplan

import (
	"io"

	"github.com/juju/ratelimit"
	"github.com/takumakei/go-timespan"
)

type Limiter struct {
	io.Writer
	timespan.Span
}

func NewLimiter(w io.Writer, bucket *ratelimit.Bucket, span timespan.Span) *Limiter {
	return &Limiter{
		Writer: ratelimit.Writer(w, bucket),
		Span:   span,
	}
}
