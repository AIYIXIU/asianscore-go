package asianscore

import (
	"math"
	"testing"
)

func TestDevigPower(t *testing.T) {
	fair, err := DevigPower([]float64{1.95, 3.50, 4.10})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	sum := fair[0] + fair[1] + fair[2]
	if math.Abs(sum-1.0) > 1e-6 {
		t.Fatalf("sum=%f want 1.0", sum)
	}
}

func TestKellyAndDixonColes(t *testing.T) {
	if got := KellyFraction(0.55, 2.0, 1.0); math.Abs(got-0.10) > 1e-6 {
		t.Fatalf("kelly=%f want 0.10", got)
	}
	p := DixonColes(1.65, 1.05, -0.05)
	if math.Abs((p.Home+p.Draw+p.Away)-1.0) > 1e-6 {
		t.Fatalf("probs do not sum to 1: %+v", p)
	}
	if p.Home <= p.Away {
		t.Fatalf("expected home > away: %+v", p)
	}
}
