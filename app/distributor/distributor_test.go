package distributor

import "testing"

func TestCalcRates(t *testing.T) {
	var calcTests = []struct {
		inN       string
		inC       string
		inCount   int
		expectedN string
		expectedC string
	}{
		{"5", "1000", 5, "1", "1"},
		{"5000", "1000", 5, "1000", "200"},
		{"10000", "1000", 5, "2000", "200"},
	}

	for _, tt := range calcTests {
		actualN, actualC, err := calcRates(tt.inN, tt.inC, tt.inCount)
		if actualN != tt.expectedN || actualC != tt.expectedC {
			t.Errorf("calcRates('%s', '%s', %d): expected ('%s', '%s'), actual ('%s','%s')", tt.inN, tt.inC, tt.inCount, tt.expectedN, tt.expectedC, actualN, actualC)
		}
		if err != nil {
			t.Errorf("calcRates(%s, %s, %d): threw error %s)", tt.inN, tt.inC, tt.inCount, err.Error())
		}
	}
}
