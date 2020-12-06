package rateplan

import (
	"fmt"
	"io"

	"github.com/CodisLabs/codis/pkg/utils/bytesize"
	"github.com/juju/ratelimit"
	"github.com/takumakei/go-timespan"
)

type RatePlan struct {
	bps    int64
	bucket *ratelimit.Bucket
	span   timespan.Span
}

func NewRatePlan(bytesPerSec int64, span timespan.Span) *RatePlan {
	return NewRatePlanClock(bytesPerSec, span, nil)
}

func NewRatePlanClock(bytesPerSec int64, span timespan.Span, clock Clock) *RatePlan {
	return &RatePlan{
		bps:    bytesPerSec,
		bucket: ratelimit.NewBucketWithRateAndClock(float64(bytesPerSec), bytesPerSec, clock),
		span:   span,
	}
}

func (rp *RatePlan) NewLimiter(w io.Writer) *Limiter {
	return NewLimiter(w, rp.bucket, rp.span)
}

func (rp *RatePlan) String() string {
	return fmt.Sprintf("%s@%s", bytesize.Int64(rp.bps).HumanString(), rp.span)
}
