package rateplan

import "io"

type Writer struct {
	w     io.Writer
	limis []*Limiter
	clock Clock
}

var _ io.Writer = (*Writer)(nil)

func NewWriter(w io.Writer, plans RatePlans) *Writer {
	return NewWriterClock(w, plans, nil)
}

func NewWriterClock(w io.Writer, plans RatePlans, clock Clock) *Writer {
	if clock == nil {
		clock = realClock{}
	}
	return &Writer{w: w, limis: plans.NewLimiters(w), clock: clock}
}

func (w *Writer) Write(b []byte) (int, error) {
	now := w.clock.Now()
	for _, limi := range w.limis {
		if limi.Contains(now) {
			return limi.Write(b)
		}
	}
	return w.w.Write(b)
}
