// Copyright 2017 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
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
