package plugin

import "testing"

func TestCompareInt64(t *testing.T) {
	type testCase struct {
		a      []int64
		b      []int64
		result bool
	}
	cases := []testCase{
		{
			a:      []int64{},
			b:      []int64{},
			result: true,
		},
		{
			a:      []int64{1, 2, 3, 4, 5},
			b:      []int64{5, 4, 3, 2, 1},
			result: true,
		},
		{
			a:      []int64{1, 2, 3, 4, 5},
			b:      []int64{1, 2, 3, 4},
			result: false,
		},
	}
	for _, testCase := range cases {
		if CompareInt64Slices(testCase.a, testCase.b) != testCase.result {
			t.Errorf("Expected %t when comparing %v == %v", testCase.result, testCase.a, testCase.b)
		}
	}
}

func TestCompareString(t *testing.T) {
	type testCase struct {
		a      []string
		b      []string
		result bool
	}
	cases := []testCase{
		{
			a:      []string{},
			b:      []string{},
			result: true,
		},
		{
			a:      []string{"1", "1", "2", "3"},
			b:      []string{"1", "1", "2", "3"},
			result: true,
		},
		{
			a:      []string{"1", "2", "3", "4", "5"},
			b:      []string{"1", "2", "3", "4"},
			result: false,
		},
	}
	for _, testCase := range cases {
		if CompareStringSlices(testCase.a, testCase.b) != testCase.result {
			t.Errorf("Expected %t when comparing %v == %v", testCase.result, testCase.a, testCase.b)
		}
	}
}

func TestCompareStringMap(t *testing.T) {
	type testCase struct {
		a      map[string]string
		b      map[string]string
		result bool
	}
	cases := []testCase{
		{
			a:      map[string]string{},
			b:      map[string]string{},
			result: true,
		},
		{
			a:      map[string]string{"1": "a", "2": "b", "3": "c", "4": "d"},
			b:      map[string]string{"1": "a", "2": "b", "3": "c", "4": "d"},
			result: true,
		},
		{
			a:      map[string]string{"1": "z", "2": "b", "3": "c", "4": "d"},
			b:      map[string]string{"1": "a", "2": "b", "3": "c", "4": "d"},
			result: false,
		},
		{
			a:      map[string]string{"1": "a", "2": "b", "3": "c", "4": "d"},
			b:      map[string]string{"1": "a", "2": "b", "3": "c"},
			result: false,
		},
		{
			a:      map[string]string{"1": "a", "2": "b"},
			b:      map[string]string{"1": "a", "2": "b", "3": "c"},
			result: false,
		},
	}
	for _, testCase := range cases {
		if CompareMapString(testCase.a, testCase.b) != testCase.result {
			t.Errorf("Expected %t when comparing %v == %v", testCase.result, testCase.a, testCase.b)
		}
	}
}

func TestCompareInt64Map(t *testing.T) {
	type testCase struct {
		a      map[string]int64
		b      map[string]int64
		result bool
	}
	cases := []testCase{
		{
			a:      map[string]int64{},
			b:      map[string]int64{},
			result: true,
		},
		{
			a:      map[string]int64{"1": 1, "2": 2, "3": 3, "4": 4},
			b:      map[string]int64{"1": 1, "2": 2, "3": 3, "4": 4},
			result: true,
		},
		{
			a:      map[string]int64{"1": 23, "2": 2, "3": 3, "4": 4},
			b:      map[string]int64{"1": 1, "2": 2, "3": 3, "4": 4},
			result: false,
		},
		{
			a:      map[string]int64{"1": 1, "2": 2, "3": 3, "4": 4},
			b:      map[string]int64{"1": 1, "2": 2, "3": 3},
			result: false,
		},
		{
			a:      map[string]int64{"1": 1, "2": 2},
			b:      map[string]int64{"1": 1, "2": 2, "3": 3},
			result: false,
		},
	}
	for _, testCase := range cases {
		if CompareMapInt64(testCase.a, testCase.b) != testCase.result {
			t.Errorf("Expected %t when comparing %v == %v", testCase.result, testCase.a, testCase.b)
		}
	}
}

func TestCompareBoolMap(t *testing.T) {
	type testCase struct {
		a      map[string]bool
		b      map[string]bool
		result bool
	}
	cases := []testCase{
		{
			a:      map[string]bool{},
			b:      map[string]bool{},
			result: true,
		},
		{
			a:      map[string]bool{"1": true, "2": false, "3": true, "4": false},
			b:      map[string]bool{"1": true, "2": false, "3": true, "4": false},
			result: true,
		},
		{
			a:      map[string]bool{"1": true, "2": false, "3": true, "4": true},
			b:      map[string]bool{"1": true, "2": false, "3": true, "4": false},
			result: false,
		},
		{
			a:      map[string]bool{"1": true, "2": true, "3": true, "4": false},
			b:      map[string]bool{"1": true, "2": true, "3": true},
			result: false,
		},
		{
			a:      map[string]bool{"1": true, "2": true},
			b:      map[string]bool{"1": true, "2": true, "3": false},
			result: false,
		},
	}
	for _, testCase := range cases {
		if CompareMapBool(testCase.a, testCase.b) != testCase.result {
			t.Errorf("Expected %t when comparing %v == %v", testCase.result, testCase.a, testCase.b)
		}
	}
}
