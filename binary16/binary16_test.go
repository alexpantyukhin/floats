package binary16

import (
	"math"
	"math/big"
	"testing"
)

func TestNewFromBits(t *testing.T) {
	golden := []struct {
		bits uint16
		want float64
	}{
		// Special numbers.

		// +NaN
		{bits: 0x7E00, want: math.NaN()},
		// -NaN
		{bits: 0xFE00, want: -math.NaN()},
		// +Inf
		{bits: 0x7C00, want: math.Inf(1)},
		// -Inf
		{bits: 0xFC00, want: math.Inf(-1)},
		// +0
		{bits: 0x0000, want: 0.0},
		// -0
		{bits: 0x8000, want: math.Copysign(0.0, -1)},

		// From https://reviews.llvm.org/rL237161

		// Normalized numbers.
		{bits: 0x3800, want: 0.5},
		{bits: 0xB800, want: -0.5},
		{bits: 0x3E00, want: 1.5},
		{bits: 0xBE00, want: -1.5},
		{bits: 0x4100, want: 2.5},
		{bits: 0xC100, want: -2.5},

		// Denormalized numbers.
		{bits: 0x0010, want: float64FromString("0x1.0p-20")},
		{bits: 0x0001, want: float64FromString("0x1.0p-24")},
		{bits: 0x8001, want: float64FromString("-0x1.0p-24")},
		//{bits: 0x0001, want: float64FromString("0x1.5p-25")},

		// Rounding.
		// TODO: Handle rounding.
		//{bits: 0x4248, want: 3.14},
		//{bits: 0xC248, want: -3.14},
		//{bits: 0x4248, want: 3.1415926535},
		//{bits: 0xC248, want: -3.1415926535},
		//{bits: 0x7C00, want: float64FromString("0x1.987124876876324p+100")},
		{bits: 0x6E62, want: float64FromString("0x1.988p+12")},
		{bits: 0x3C00, want: float64FromString("0x1.0p+0")},
		{bits: 0x0400, want: float64FromString("0x1.0p-14")},
		// rounded to zero
		//{bits: 0x0000, want: float64FromString("0x1.0p-25")},
		//{bits: 0x8000, want: float64FromString("-0x1.0p-25")},
		// max (precise)
		{bits: 0x7BFF, want: 65504.0},
	}

	for _, g := range golden {
		f := NewFromBits(g.bits)
		got := f.Float64()
		wantBits := math.Float64bits(g.want)
		gotBits := math.Float64bits(got)
		if wantBits != gotBits {
			t.Errorf("0x%04X: number mismatch; expected 0x%08X (%v), got 0x%08X (%v)", g.bits, wantBits, g.want, gotBits, got)
		}
	}
}

func float64FromString(s string) float64 {
	x, _, err := big.ParseFloat(s, 0, 53, big.ToNearestEven)
	if err != nil {
		panic(err)
	}
	// TODO: Check accuracy?
	y, _ := x.Float64()
	return y
}
