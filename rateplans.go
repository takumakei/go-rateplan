package rateplan

import "io"

type RatePlans []*RatePlan

func (rps RatePlans) NewLimiters(w io.Writer) []*Limiter {
	limis := make([]*Limiter, len(rps))
	for i, rp := range rps {
		limis[i] = rp.NewLimiter(w)
	}
	return limis
}
