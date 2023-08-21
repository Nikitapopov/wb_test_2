package main

import (
	"reflect"
	"testing"
)

type testCase struct {
	input  *[]string
	gotRes map[string]*[]string
}

func test(tc testCase, t *testing.T) {
	res := getAnagrams(tc.input)
	if !reflect.DeepEqual(res, tc.gotRes) {
		t.Errorf("expected result: %v, got: %v", tc.gotRes, res)
	}
}

func Test_getAnagrams_cases_order_unique_1(t *testing.T) {
	tc := testCase{
		input: &[]string{"пятак", "пятка", "кот", "лиСток", "тяпка", "столик", "пятак", "слиток"},
		gotRes: map[string]*[]string{
			"пятак":  {"пятка", "тяпка"},
			"листок": {"слиток", "столик"},
		},
	}
	test(tc, t)
}

func Test_getAnagrams_cases_order_unique_2(t *testing.T) {
	tc := testCase{
		input: &[]string{"собака", "пятка", "столик", "пятак", "ЛИСТОК", "тяпка", "слиток", "слиток", "ислток"},
		gotRes: map[string]*[]string{
			"пятка":  {"пятак", "тяпка"},
			"столик": {"ислток", "листок", "слиток"},
		},
	}
	test(tc, t)
}

func Test_getAnagrams_empty_list(t *testing.T) {
	tc := testCase{
		input:  &[]string{},
		gotRes: map[string]*[]string{},
	}
	test(tc, t)
}

func Test_getAnagrams_same_word(t *testing.T) {
	tc := testCase{
		input:  &[]string{"собака", "собака", "собака"},
		gotRes: map[string]*[]string{},
	}
	test(tc, t)
}
