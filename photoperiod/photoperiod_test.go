package photoperiod

import "testing"

func Test_percentToIntensity(t *testing.T) {
	tests := []struct {
		p float64
		want int
	}{
		{0,0},
		{0.01, 10},
		{.1,100},
		{.5,500},
		{1,1000},
		{1.25,2000},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if got := percentToIntensity(tt.p); got != tt.want {
				t.Errorf("percentToIntensity() = %v, want %v", got, tt.want)
			}
		})
	}
}
