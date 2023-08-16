package main

import (
	"testing"
)

type testCase struct {
	input  string
	gotRes string
	gotErr string
}

func test(tc testCase, t *testing.T) {
	res, err := parseStr(tc.input)
	if tc.gotErr != "" {
		if err == nil {
			t.Errorf("expected error: %s", tc.gotErr)
		} else if tc.gotErr != err.Error() {
			t.Errorf("expected error: %s, got: %s", tc.gotErr, err.Error())
		}
	} else {
		if res != tc.gotRes {
			t.Errorf("expected result: %s, got: %s", tc.gotRes, res)
		}
	}
}

func Test_parseStr_abc3d12e_escape_5f(t *testing.T) {
	tc := testCase{
		input:  `abc3d12e\5f`,
		gotRes: "abcccdddddddddddde5f",
	}
	test(tc, t)
}

func Test_parseStr_a4bc2d5e(t *testing.T) {
	tc := testCase{
		input:  "a4bc2d5e",
		gotRes: "aaaabccddddde",
	}
	test(tc, t)
}

func Test_parseStr_abcd(t *testing.T) {
	tc := testCase{
		input:  "abcd",
		gotRes: "abcd",
	}
	test(tc, t)
}

func Test_parseStr_empty_string(t *testing.T) {
	tc := testCase{
		input:  "",
		gotRes: "",
	}
	test(tc, t)
}

func Test_parseStr_a11(t *testing.T) {
	tc := testCase{
		input:  "a11",
		gotRes: "aaaaaaaaaaa",
	}
	test(tc, t)
}

func Test_parseStr_a_escape_11(t *testing.T) {
	tc := testCase{
		input:  `a\11`,
		gotRes: "a1",
	}
	test(tc, t)
}

func Test_parseStr_a1_escape_1(t *testing.T) {
	tc := testCase{
		input:  `a1\1`,
		gotRes: "a1",
	}
	test(tc, t)
}

func Test_parseStr_qwe_escape_4_escape_5(t *testing.T) {
	tc := testCase{
		input:  `qwe\4\5`,
		gotRes: "qwe45",
	}
	test(tc, t)
}

func Test_parseStr_qwe_escape_45(t *testing.T) {
	tc := testCase{
		input:  `qwe\45`,
		gotRes: "qwe44444",
	}
	test(tc, t)
}

func Test_parseStr_qwe_escape_escape_5(t *testing.T) {
	tc := testCase{
		input:  `qwe\\5`,
		gotRes: `qwe\\\\\`,
	}
	test(tc, t)
}

func Test_parseStr_a03(t *testing.T) {
	tc := testCase{
		input:  `a03`,
		gotRes: "aaa",
	}
	test(tc, t)
}

func Test_parseStr_a_minus_2(t *testing.T) {
	tc := testCase{
		input:  `a-2`,
		gotRes: "a--",
	}
	test(tc, t)
}

func Test_parseStr_45(t *testing.T) {
	tc := testCase{
		input:  `45`,
		gotErr: "некорректная строка",
	}
	test(tc, t)
}

func Test_parseStr_1(t *testing.T) {
	tc := testCase{
		input:  "1",
		gotErr: "некорректная строка",
	}
	test(tc, t)
}

func Test_parseStr_escape(t *testing.T) {
	tc := testCase{
		input:  `\`,
		gotErr: "некорректная строка",
	}
	test(tc, t)
}
