package rateplan

import (
	"errors"
	"fmt"
	"strings"

	"github.com/inhies/go-bytesize"
	"github.com/takumakei/go-timespan"
)

var (
	ErrParseRatePlan = errors.New("can't convert a string into RatePlan")
)

func ParseRatePlans(s string) (RatePlans, error) {
	a := strings.Split(s, ",")
	plans := make([]*RatePlan, len(a))
	for i, e := range a {
		plan, err := ParseRatePlan(e)
		if err != nil {
			return nil, err
		}
		plans[i] = plan
	}
	return plans, nil
}

func ParseRatePlan(s string) (*RatePlan, error) {
	s = strings.TrimSpace(s)
	a := strings.SplitN(s, "@", 2)
	if len(a) != 2 {
		return nil, fmt.Errorf("%w: %q", ErrParseRatePlan, s)
	}
	bytesPerSec, err := bytesize.Parse(a[0])
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrParseRatePlan, fmt.Errorf("%w: %q", err, a[0]))
	}
	span, err := timespan.Parse(strings.TrimSpace(a[1]))
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrParseRatePlan, fmt.Errorf("%w, %q", err, a[1]))
	}
	return NewRatePlan(int64(bytesPerSec), span), nil
}
