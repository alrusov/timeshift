package timeshift

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/alrusov/misc"
)

//----------------------------------------------------------------------------------------------------------------------------//

var testParameters = []struct {
	pattern       string
	errorExpected bool
	t             time.Time
	result        time.Time
}{
	{pattern: "ZZZYYY", errorExpected: true},
	{pattern: "Z1M2D$-3W3w1h-4m+5s6", errorExpected: true},
	{pattern: "YM2D$-3h-4m+5s6", errorExpected: true},
	{pattern: "Y1M2D^$-3h-4m+5s6", errorExpected: true},
	{pattern: "Y1M2D$*3h-4m+5s6", errorExpected: true},
	{pattern: "Y1M2D$-h-4m+5s6", errorExpected: true},
	{pattern: "Y1M2D$-3hm+5s6", errorExpected: true},
	{pattern: "Y1 M2 D$-3 H-4 m+5 s6", errorExpected: true},
	{pattern: "Y1M2D-3h$-4m+5s6", errorExpected: true},
	{pattern: "  Y$-2  ", errorExpected: true},
	{pattern: "Y+1 M+2 D$+3 h-6 m+20 s-30", errorExpected: true},
	{pattern: "Y1 M2 D$-3 h-4 h+5 s6", errorExpected: true},
	{pattern: "M1 Y2 D$-3 h-4 m+5 s6", errorExpected: true},
	{pattern: "Y+1 M+2 D^-3 h-6 m+20 s-30", errorExpected: true},
	{pattern: "D0", errorExpected: true},
	{pattern: "M0", errorExpected: true},
	{pattern: "W0", errorExpected: true},
	{pattern: "W^+0", errorExpected: true},
	{pattern: "W$0", errorExpected: true},
	{pattern: "W^-2", errorExpected: true},
	{pattern: "W$+2", errorExpected: true},
	{pattern: "W$^2", errorExpected: true},
	{pattern: "w-1", errorExpected: true},
	{pattern: "w8", errorExpected: true},
	{pattern: "w$+3", errorExpected: true},
	{pattern: "w^-3", errorExpected: true},

	{pattern: "", errorExpected: false, t: tConv("2020-06-13T14:55:22Z"), result: tConv("2020-06-13T14:55:22Z")},
	{pattern: "", errorExpected: false, t: tConv("2020-06-13T14:55:21+03:00"), result: tConv("2020-06-13T14:55:21+03:00")},
	{pattern: "    ", errorExpected: false, t: tConv("2020-06-13T11:55:22+03:00"), result: tConv("2020-06-13T11:55:22+03:00")},
	{pattern: "  Y2021  ", errorExpected: false, t: tConv("2020-06-13T14:55:22Z"), result: tConv("2021-06-13T14:55:22Z")},
	{pattern: "  Y2021 M2 D3 h6 m20 s30", errorExpected: false, t: tConv("2020-06-13T14:55:22Z"), result: tConv("2021-02-03T06:20:30Z")},
	{pattern: "  Y2021 M22 D33 h66 m200 s300", errorExpected: false, t: tConv("2020-06-13T14:55:22Z"), result: tConv("2022-11-04T21:25:00Z")},

	{pattern: "  Y+1 M+22 D+33 h+66 m+200 s+300", errorExpected: false, t: tConv("2020-06-13T14:55:22Z"), result: tConv("2023-05-19T12:20:22Z")},
	{pattern: "Y+1M+22D+33h+66m+200s+300", errorExpected: false, t: tConv("2020-06-13T14:55:22Z"), result: tConv("2023-05-19T12:20:22Z")},
	{pattern: "  Y+1        M+22D+33h+66      m+200           s+300      ", errorExpected: false, t: tConv("2020-06-13T14:55:22Z"), result: tConv("2023-05-19T12:20:22Z")},
	{pattern: "  Y+1 M+2 D$3 h-6 m+20 s-30", errorExpected: false, t: tConv("2020-06-13T14:55:22Z"), result: tConv("2021-08-29T09:14:52Z")},
	{pattern: "  Y+1 M+2 D$3 W+2 h-6 m+20 s-30", errorExpected: false, t: tConv("2020-06-13T14:55:22Z"), result: tConv("2021-09-12T09:14:52Z")},
	{pattern: "  Y+1 M+2 D$3 W-2 h-6 m+20 s-30", errorExpected: false, t: tConv("2020-06-13T14:55:22Z"), result: tConv("2021-08-15T09:14:52Z")},
	{pattern: "D$1", errorExpected: false, t: tConv("2000-02-13T14:55:22Z"), result: tConv("2000-02-29T14:55:22Z")},
	{pattern: "Y+100 D$1", errorExpected: false, t: tConv("2000-02-13T14:55:22Z"), result: tConv("2100-02-28T14:55:22Z")},
	{pattern: "Y-100 D$1", errorExpected: false, t: tConv("2000-02-13T14:55:22Z"), result: tConv("1900-02-28T14:55:22Z")},
	{pattern: "  Y2021 M-3 D-13 h-6 m-56 s-23", errorExpected: false, t: tConv("2020-02-13T14:55:22Z"), result: tConv("2020-10-31T07:58:59Z")},
	{pattern: "  Y2021 M+13 D+20 h-6 m-56 s-23", errorExpected: false, t: tConv("2020-02-13T14:55:22Z"), result: tConv("2022-04-02T07:58:59Z")},
	{pattern: "  Y2021 M13 D33 h47 m62 s125", errorExpected: false, t: tConv("2020-02-13T14:55:22Z"), result: tConv("2022-02-04T00:04:05Z")},

	{pattern: "w0", errorExpected: false, t: tConv("2021-02-22T14:55:22Z"), result: tConv("2021-02-21T14:55:22Z")},
	{pattern: "w1", errorExpected: false, t: tConv("2021-02-22T14:55:22Z"), result: tConv("2021-02-22T14:55:22Z")},
	{pattern: "w2", errorExpected: false, t: tConv("2021-02-22T14:55:22Z"), result: tConv("2021-02-23T14:55:22Z")},
	{pattern: "w3", errorExpected: false, t: tConv("2021-02-22T14:55:22Z"), result: tConv("2021-02-24T14:55:22Z")},
	{pattern: "w4", errorExpected: false, t: tConv("2021-02-22T14:55:22Z"), result: tConv("2021-02-25T14:55:22Z")},
	{pattern: "w5", errorExpected: false, t: tConv("2021-02-22T14:55:22Z"), result: tConv("2021-02-26T14:55:22Z")},
	{pattern: "w6", errorExpected: false, t: tConv("2021-02-22T14:55:22Z"), result: tConv("2021-02-27T14:55:22Z")},

	{pattern: "w0", errorExpected: false, t: tConv("2021-02-07T14:55:22Z"), result: tConv("2021-02-07T14:55:22Z")},
	{pattern: "w1", errorExpected: false, t: tConv("2021-02-07T14:55:22Z"), result: tConv("2021-02-08T14:55:22Z")},

	{pattern: "w0 h2 m0 s0", errorExpected: false, t: tConv("2021-02-22T14:55:22Z"), result: tConv("2021-02-21T02:00:00Z")},
	{pattern: "w1 h2 m0 s0", errorExpected: false, t: tConv("2021-02-22T14:55:22Z"), result: tConv("2021-02-22T02:00:00Z")},
	{pattern: "w5 h2 m0 s0", errorExpected: false, t: tConv("2021-02-22T14:55:22Z"), result: tConv("2021-02-26T02:00:00Z")},

	{pattern: "W1 w2", errorExpected: false, t: tConv("2020-01-01T00:00:00Z"), result: tConv("2020-01-07T00:00:00Z")},
	{pattern: "W2 w2", errorExpected: false, t: tConv("2020-01-01T00:00:00Z"), result: tConv("2020-01-14T00:00:00Z")},
	{pattern: "W51 w2", errorExpected: false, t: tConv("2020-01-01T00:00:00Z"), result: tConv("2020-12-22T00:00:00Z")},
	{pattern: "W52 w2", errorExpected: false, t: tConv("2020-10-01T00:00:00Z"), result: tConv("2020-12-29T00:00:00Z")},
	{pattern: "W53 w2", errorExpected: false, t: tConv("2020-01-01T00:00:00Z"), result: tConv("2021-01-05T00:00:00Z")},
	{pattern: "W54 w2", errorExpected: false, t: tConv("2020-01-01T00:00:00Z"), result: tConv("2021-01-12T00:00:00Z")},

	{pattern: "W1 w3", errorExpected: false, t: tConv("2020-01-01T00:00:00Z"), result: tConv("2020-01-01T00:00:00Z")},
	{pattern: "W2 w3", errorExpected: false, t: tConv("2020-01-01T00:00:00Z"), result: tConv("2020-01-08T00:00:00Z")},

	{pattern: "W1 w5", errorExpected: false, t: tConv("2021-01-02T00:00:00Z"), result: tConv("2021-01-01T00:00:00Z")},
	{pattern: "W1 w6", errorExpected: false, t: tConv("2021-01-02T00:00:00Z"), result: tConv("2021-01-02T00:00:00Z")},
	{pattern: "W1 w1", errorExpected: false, t: tConv("2021-01-02T00:00:00Z"), result: tConv("2021-01-04T00:00:00Z")},
	{pattern: "W1 w2", errorExpected: false, t: tConv("2021-01-02T00:00:00Z"), result: tConv("2021-01-05T00:00:00Z")},
	{pattern: "W1 w3", errorExpected: false, t: tConv("2021-01-02T00:00:00Z"), result: tConv("2021-01-06T00:00:00Z")},
	{pattern: "W1 w4", errorExpected: false, t: tConv("2021-01-02T00:00:00Z"), result: tConv("2021-01-07T00:00:00Z")},
	{pattern: "W1 w0", errorExpected: false, t: tConv("2021-01-02T00:00:00Z"), result: tConv("2021-01-03T00:00:00Z")},

	{pattern: "W2 w2", errorExpected: false, t: tConv("2021-01-01T00:00:00Z"), result: tConv("2021-01-12T00:00:00Z")},
	{pattern: "W1 w2", errorExpected: false, t: tConv("2021-11-08T00:00:00Z"), result: tConv("2021-01-05T00:00:00Z")},
	{pattern: "W2 w2", errorExpected: false, t: tConv("2021-10-08T00:00:00Z"), result: tConv("2021-01-12T00:00:00Z")},

	{pattern: "W1 w3", errorExpected: false, t: tConv("2020-11-01T00:00:00Z"), result: tConv("2020-01-01T00:00:00Z")},
	{pattern: "W1 w4", errorExpected: false, t: tConv("2020-11-01T00:00:00Z"), result: tConv("2020-01-02T00:00:00Z")},
	{pattern: "W1 w5", errorExpected: false, t: tConv("2020-11-01T00:00:00Z"), result: tConv("2020-01-03T00:00:00Z")},
	{pattern: "W1 w6", errorExpected: false, t: tConv("2020-11-01T00:00:00Z"), result: tConv("2020-01-04T00:00:00Z")},
	{pattern: "W1 w1", errorExpected: false, t: tConv("2020-11-01T00:00:00Z"), result: tConv("2020-01-06T00:00:00Z")},
	{pattern: "W1 w2", errorExpected: false, t: tConv("2020-11-01T00:00:00Z"), result: tConv("2020-01-07T00:00:00Z")},
	{pattern: "W1 w0", errorExpected: false, t: tConv("2020-11-01T00:00:00Z"), result: tConv("2020-01-05T00:00:00Z")},

	{pattern: "W1 w0", errorExpected: false, t: tConv("2017-11-01T00:00:00Z"), result: tConv("2017-01-01T00:00:00Z")},

	{pattern: "W1 w1", errorExpected: false, t: tConv("2017-11-01T00:00:00Z"), result: tConv("2017-01-02T00:00:00Z")},

	{pattern: "W51 w0", errorExpected: false, t: tConv("2020-01-01T00:00:00Z"), result: tConv("2020-12-20T00:00:00Z")},
	{pattern: "W52 w0", errorExpected: false, t: tConv("2020-10-01T00:00:00Z"), result: tConv("2020-12-27T00:00:00Z")},
	{pattern: "W53 w0", errorExpected: false, t: tConv("2020-01-01T00:00:00Z"), result: tConv("2021-01-03T00:00:00Z")},
	{pattern: "W54 w0", errorExpected: false, t: tConv("2020-01-01T00:00:00Z"), result: tConv("2021-01-10T00:00:00Z")},
	{pattern: "W1 w0", errorExpected: false, t: tConv("2021-01-01T00:00:00Z"), result: tConv("2021-01-03T00:00:00Z")},
	{pattern: "W1 w0", errorExpected: false, t: tConv("2021-12-08T00:00:00Z"), result: tConv("2021-01-03T00:00:00Z")},

	{pattern: "W-1 w2", errorExpected: false, t: tConv("2021-02-01T01:00:00Z"), result: tConv("2021-01-26T01:00:00Z")},
	{pattern: "W+0 w2", errorExpected: false, t: tConv("2021-02-01T02:00:00Z"), result: tConv("2021-02-02T02:00:00Z")},
	{pattern: "W-0 w2", errorExpected: false, t: tConv("2021-02-01T03:00:00Z"), result: tConv("2021-02-02T03:00:00Z")},
	{pattern: "W+1 w2", errorExpected: false, t: tConv("2021-02-01T04:00:00Z"), result: tConv("2021-02-09T04:00:00Z")},

	{pattern: "D+6 W-1 w2", errorExpected: false, t: tConv("2021-02-01T00:11:00Z"), result: tConv("2021-02-02T00:11:00Z")},
	{pattern: "D+6 W+0 w2", errorExpected: false, t: tConv("2021-02-01T00:22:00Z"), result: tConv("2021-02-09T00:22:00Z")},
	{pattern: "D+6 W-0 w2", errorExpected: false, t: tConv("2021-02-01T00:33:00Z"), result: tConv("2021-02-09T00:33:00Z")},
	{pattern: "D+6 W+1 w2", errorExpected: false, t: tConv("2021-02-01T00:44:00Z"), result: tConv("2021-02-16T00:44:00Z")},

	{pattern: "W^1 w0", errorExpected: false, t: tConv("2021-01-20T00:00:00Z"), result: tConv("2021-01-03T00:00:00Z")},
	{pattern: "W^2 w0", errorExpected: false, t: tConv("2021-01-20T00:00:00Z"), result: tConv("2021-01-10T00:00:00Z")},
	{pattern: "W^1 w5", errorExpected: false, t: tConv("2021-01-20T00:00:00Z"), result: tConv("2021-01-01T00:00:00Z")},
	{pattern: "W^10 w5", errorExpected: false, t: tConv("2021-01-20T00:00:00Z"), result: tConv("2021-03-05T00:00:00Z")},
	{pattern: "W^1 w4", errorExpected: false, t: tConv("2021-01-20T00:00:00Z"), result: tConv("2021-01-07T00:00:00Z")},

	{pattern: "W^1 w0", errorExpected: false, t: tConv("2021-08-10T00:00:00Z"), result: tConv("2021-08-01T00:00:00Z")},
	{pattern: "W^1 w0", errorExpected: false, t: tConv("2021-03-20T00:00:00Z"), result: tConv("2021-03-07T00:00:00Z")},

	{pattern: "W$1 w0", errorExpected: false, t: tConv("2021-03-20T00:00:00Z"), result: tConv("2021-03-28T00:00:00Z")},
	{pattern: "W$1 w1", errorExpected: false, t: tConv("2021-03-20T00:00:00Z"), result: tConv("2021-03-29T00:00:00Z")},
	{pattern: "W$1 w2", errorExpected: false, t: tConv("2021-03-20T00:00:00Z"), result: tConv("2021-03-30T00:00:00Z")},
	{pattern: "W$1 w3", errorExpected: false, t: tConv("2021-03-20T00:00:00Z"), result: tConv("2021-03-31T00:00:00Z")},
	{pattern: "W$1 w4", errorExpected: false, t: tConv("2021-03-20T00:00:00Z"), result: tConv("2021-03-25T00:00:00Z")},
	{pattern: "W$1 w5", errorExpected: false, t: tConv("2021-03-20T00:00:00Z"), result: tConv("2021-03-26T00:00:00Z")},
	{pattern: "W$1 w6", errorExpected: false, t: tConv("2021-03-20T00:00:00Z"), result: tConv("2021-03-27T00:00:00Z")},

	{pattern: "W$4 w0", errorExpected: false, t: tConv("2021-03-20T00:00:00Z"), result: tConv("2021-03-07T00:00:00Z")},
	{pattern: "W$4 w1", errorExpected: false, t: tConv("2021-03-20T00:00:00Z"), result: tConv("2021-03-08T00:00:00Z")},
	{pattern: "W$4 w2", errorExpected: false, t: tConv("2021-03-20T00:00:00Z"), result: tConv("2021-03-09T00:00:00Z")},
	{pattern: "W$4 w3", errorExpected: false, t: tConv("2021-03-20T00:00:00Z"), result: tConv("2021-03-10T00:00:00Z")},
	{pattern: "W$4 w4", errorExpected: false, t: tConv("2021-03-20T00:00:00Z"), result: tConv("2021-03-04T00:00:00Z")},
	{pattern: "W$4 w5", errorExpected: false, t: tConv("2021-03-20T00:00:00Z"), result: tConv("2021-03-05T00:00:00Z")},
	{pattern: "W$4 w6", errorExpected: false, t: tConv("2021-03-20T00:00:00Z"), result: tConv("2021-03-06T00:00:00Z")},

	{pattern: "l+10 u-2 n+1234", errorExpected: false, t: tConv("2021-03-20T00:00:00Z"), result: tConv("2021-03-20T00:00:00.009999234Z")},
}

func tConv(s string) time.Time {
	tt, err := misc.ParseJSONtime(s)
	if err != nil {
		log.Fatal(err)
	}
	return tt
}

//----------------------------------------------------------------------------------------------------------------------------//

func Test1(t *testing.T) {
	for i, p := range testParameters {
		ts, err := New(p.pattern, false)

		if p.errorExpected {
			if err == nil {
				t.Errorf(`[%d] "%s" prepared without error, expected error`, i, p.pattern)
			}
			continue
		}

		if err != nil {
			t.Errorf(`[%d] "%s" prepared with error: %s`, i, p.pattern, err)
			continue
		}

		result := ts.Exec(p.t)

		//if !result.Equal(p.result) {
		if result != p.result { // the locale must be saved!
			t.Errorf(`[%d] "%s" shifted by "%s": got "%s", expected "%s"`, i, misc.Time2JSONtz(p.t), p.pattern, misc.Time2JSONtz(result), misc.Time2JSONtz(p.result))
		}
	}
}

//----------------------------------------------------------------------------------------------------------------------------//

func BenchmarkCacheOn(b *testing.B) {
	benchmark(b, true)
}

func BenchmarkCacheOff(b *testing.B) {
	benchmark(b, false)
}

func benchmark(b *testing.B, cached bool) {
	pattern := "Y+1 M+2 D$3 W-2 h-6 m+20 s-30"
	t := tConv("2020-06-13T14:55:22Z")
	expected := tConv("2021-08-15T09:14:52Z")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ts, err := New(pattern, cached)
		if err != nil {
			b.Fatalf(`[%d] "%s" prepared with error: %s`, i, pattern, err)
		}

		result := ts.Exec(t)

		if result != expected { // the locale must be saved!
			b.Fatalf(`[%d] "%s" shifted by "%s": got "%s", expected "%s"`, i, misc.Time2JSONtz(t), pattern, misc.Time2JSONtz(result), misc.Time2JSONtz(expected))
		}
	}
}

//----------------------------------------------------------------------------------------------------------------------------//

func TestPrintParameters(t *testing.T) {
	for _, p := range testParameters {
		if !p.errorExpected {
			fmt.Printf("|\"%s\"|%s|%s|\n", p.pattern, misc.Time2JSONtz(p.t), misc.Time2JSONtz(p.result))
		}
	}
}

//----------------------------------------------------------------------------------------------------------------------------//
