package rateplan_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/takumakei/go-rateplan"
)

func TestParseRatePlans(t *testing.T) {
	t.Run("123", func(t *testing.T) {
		_, err := rateplan.ParseRatePlan("123")
		if err == nil {
			t.Error(`rateplan.ParseRatePlan("123")`)
		}
		//t.Log(err)
	})

	t.Run("123x@", func(t *testing.T) {
		_, err := rateplan.ParseRatePlan("123x@")
		if err == nil {
			t.Error(`rateplan.ParseRatePlan("123x@")`)
		}
		//t.Log(err)
	})

	t.Run("128kb@abc", func(t *testing.T) {
		_, err := rateplan.ParseRatePlan("128kb@abc/def")
		if err == nil {
			t.Error(`rateplan.ParseRatePlan("128kb@abc/def")`)
		}
		//t.Log(err)
	})

	t.Run("ok", func(t *testing.T) {
		a, err := rateplan.ParseRatePlans("1mb@02:00+0900/2h,10mb@24h")
		if err != nil {
			t.Error(`rateplan.ParseRatePlans("1mb@02:00+0900/2h,10mb@24h")`)
		}
		w := `[1024.00kb@02:00+0900/2h0m0s 10.00mb@24h0m0s]`
		if s := fmt.Sprint(a); s != w {
			t.Error(w)
		}
	})
}

func Test(t *testing.T) {
	t.Run("strings.SplitN", func(t *testing.T) {
		if a := strings.SplitN("", "@", 2); len(a) != 1 {
			t.Error(`a := strings.SplitN("", "@", 2); len(a) != 1`)
		}
	})
}
