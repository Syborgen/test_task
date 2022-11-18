package datastructures

import "testing"

func TestIsOverlapped(t *testing.T) {
	testCases := []struct {
		range1            TimeRange
		range2            TimeRange
		isOverlapExpected bool
	}{
		{
			range1:            TimeRange{Start: "2022-01-01 00:00:20", End: "2022-01-01 00:00:40"},
			range2:            TimeRange{Start: "2022-01-01 00:00:00", End: "2022-01-01 00:00:20"},
			isOverlapExpected: false,
		},
		{
			range1:            TimeRange{Start: "2022-01-01 00:00:00", End: "2022-01-01 00:00:20"},
			range2:            TimeRange{Start: "2022-01-01 00:00:20", End: "2022-01-01 00:00:40"},
			isOverlapExpected: false,
		},
		{
			range1:            TimeRange{Start: "2022-01-01 00:00:10", End: "2022-01-01 00:00:20"},
			range2:            TimeRange{Start: "2022-01-01 00:00:00", End: "2022-01-01 00:00:30"},
			isOverlapExpected: true,
		},
	}

	for _, tc := range testCases {
		if tc.range1.IsOverlapped(tc.range2) != tc.isOverlapExpected {
			t.Errorf("Is overlap expected: %v, actual overlap: %v",
				tc.isOverlapExpected, tc.range1.IsOverlapped(tc.range2))
		}
	}
}
