package main

import (
	"reflect"
	"testing"
)

type filterParams struct {
	lines []string
	args  cmdArgs
}

type testCase struct {
	input     filterParams
	expect    []string
	expectErr bool
}

func test(t *testing.T, tc testCase) {
	got, err := filter(tc.input.lines, tc.input.args)
	if err != nil && !tc.expectErr {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	if err != nil && !tc.expectErr {
		t.Errorf("Expected error: %v", tc.expectErr)
		return
	}

	if !reflect.DeepEqual(got, tc.expect) {
		t.Errorf("Expected: %v, got: %v", tc.expect, got)
	}
}

func Test_filter_by_default_options_with_2_matches(t *testing.T) {
	tc := testCase{
		input: filterParams{
			lines: []string{
				"Simon Strange",
				"Pete Brown",
				"Tony Bedford",
				"Stefan Heinz",
				"Tony Soprano",
				"John Strange",
				"Fred Bloggs",
			},
			args: cmdArgs{
				after:      0,
				before:     0,
				count:      false,
				ignoreCase: false,
				invert:     false,
				fixed:      false,
				lineNum:    false,
				pattern:    "Tony",
			},
		},
		expect: []string{
			"Tony Bedford",
			"Tony Soprano",
		},
		expectErr: false,
	}
	t.Run("filter by default options with 2 matches", func(t *testing.T) { test(t, tc) })
}

func Test_filter_by_default_options_with_0_matches(t *testing.T) {
	tc := testCase{
		input: filterParams{
			lines: []string{
				"Simon Strange",
				"Pete Brown",
				"Tony Bedford",
				"Stefan Heinz",
				"Tony Soprano",
				"John Strange",
				"Fred Bloggs",
			},
			args: cmdArgs{
				after:      0,
				before:     0,
				count:      false,
				ignoreCase: false,
				invert:     false,
				fixed:      false,
				lineNum:    false,
				pattern:    "Alex",
			},
		},
		expect:    []string{},
		expectErr: false,
	}
	t.Run("filter by default options with 0 matches", func(t *testing.T) { test(t, tc) })
}

func Test_filter_with_A_equal_2_flag(t *testing.T) {
	tc := testCase{
		input: filterParams{
			lines: []string{
				"Simon Strange",
				"Pete Brown",
				"Tony Bedford",
				"Stefan Heinz",
				"Tony Soprano",
				"John Strange",
				"Fred Bloggs",
			},
			args: cmdArgs{
				after:      2,
				before:     0,
				count:      false,
				ignoreCase: false,
				invert:     false,
				fixed:      false,
				lineNum:    false,
				pattern:    "Heinz",
			},
		},
		expect: []string{
			"Stefan Heinz",
			"Tony Soprano",
			"John Strange",
		},
		expectErr: false,
	}
	t.Run("filter with A=2 flag", func(t *testing.T) { test(t, tc) })
}

func Test_filter_with_A_equal_1_flag_last_element(t *testing.T) {
	tc := testCase{
		input: filterParams{
			lines: []string{
				"Simon Strange",
				"Pete Brown",
				"Tony Bedford",
				"Stefan Heinz",
				"Tony Soprano",
				"John Strange",
				"Fred Bloggs",
			},
			args: cmdArgs{
				after:      1,
				before:     0,
				count:      false,
				ignoreCase: false,
				invert:     false,
				fixed:      false,
				lineNum:    false,
				pattern:    "Fred",
			},
		},
		expect: []string{
			"Fred Bloggs",
		},
		expectErr: false,
	}
	t.Run("filter with A=2 flag", func(t *testing.T) { test(t, tc) })
}

func Test_filter_with_B_equal_2_flag(t *testing.T) {
	tc := testCase{
		input: filterParams{
			lines: []string{
				"Simon Strange",
				"Pete Brown",
				"Tony Bedford",
				"Stefan Heinz",
				"Tony Soprano",
				"John Strange",
				"Fred Bloggs",
			},
			args: cmdArgs{
				after:      0,
				before:     2,
				count:      false,
				ignoreCase: false,
				invert:     false,
				fixed:      false,
				lineNum:    false,
				pattern:    "Heinz",
			},
		},
		expect: []string{
			"Pete Brown",
			"Tony Bedford",
			"Stefan Heinz",
		},
		expectErr: false,
	}
	t.Run("filter with B=2 flag", func(t *testing.T) { test(t, tc) })
}

func Test_filter_with_B_equal_2_flag_first_element(t *testing.T) {
	tc := testCase{
		input: filterParams{
			lines: []string{
				"Simon Strange",
				"Pete Brown",
				"Tony Bedford",
				"Stefan Heinz",
				"Tony Soprano",
				"John Strange",
				"Fred Bloggs",
			},
			args: cmdArgs{
				after:      0,
				before:     1,
				count:      false,
				ignoreCase: false,
				invert:     false,
				fixed:      false,
				lineNum:    false,
				pattern:    "Simon",
			},
		},
		expect: []string{
			"Simon Strange",
		},
		expectErr: false,
	}
	t.Run("filter with B=2 flag", func(t *testing.T) { test(t, tc) })
}

func Test_filter_with_intersection_of_A_equal_1_B_equal_2_flags(t *testing.T) {
	tc := testCase{
		input: filterParams{
			lines: []string{
				"Simon Strange",
				"Pete Brown",
				"Tony Bedford",
				"Stefan Heinz",
				"Tony Soprano",
				"John Strange",
				"Fred Bloggs",
			},
			args: cmdArgs{
				after:      1,
				before:     2,
				count:      false,
				ignoreCase: false,
				invert:     false,
				fixed:      false,
				lineNum:    false,
				pattern:    "Tony",
			},
		},
		expect: []string{
			"Simon Strange",
			"Pete Brown",
			"Tony Bedford",
			"Stefan Heinz",
			"Tony Soprano",
			"John Strange",
		},
		expectErr: false,
	}
	t.Run("filter with intersection of A=1, B=2 flags", func(t *testing.T) { test(t, tc) })
}

func Test_filter_with_c_equal_true_flag(t *testing.T) {
	tc := testCase{
		input: filterParams{
			lines: []string{
				"Simon Strange",
				"Pete Brown",
				"Tony Bedford",
				"Stefan Heinz",
				"Tony Soprano",
				"John Strange",
				"Fred Bloggs",
			},
			args: cmdArgs{
				after:      0,
				before:     0,
				count:      true,
				ignoreCase: false,
				invert:     false,
				fixed:      false,
				lineNum:    false,
				pattern:    "Tony",
			},
		},
		expect: []string{
			"2",
		},
		expectErr: false,
	}
	t.Run("filter with c=true flag", func(t *testing.T) { test(t, tc) })
}

func Test_filter_with_i_equal_true_flag(t *testing.T) {
	tc := testCase{
		input: filterParams{
			lines: []string{
				"Simon Strange",
				"Pete Brown",
				"Tony Bedford",
				"Stefan Heinz",
				"Tony Soprano",
				"John Strange",
				"Fred Bloggs",
			},
			args: cmdArgs{
				after:      0,
				before:     0,
				count:      false,
				ignoreCase: true,
				invert:     false,
				fixed:      false,
				lineNum:    false,
				pattern:    "tonY",
			},
		},
		expect: []string{
			"Tony Bedford",
			"Tony Soprano",
		},
		expectErr: false,
	}
	t.Run("filter with i=true flag", func(t *testing.T) { test(t, tc) })
}

func Test_filter_with_v_equal_true_flag(t *testing.T) {
	tc := testCase{
		input: filterParams{
			lines: []string{
				"Simon Strange",
				"Pete Brown",
				"Tony Bedford",
				"Stefan Heinz",
				"Tony Soprano",
				"John Strange",
				"Fred Bloggs",
			},
			args: cmdArgs{
				after:      0,
				before:     0,
				count:      false,
				ignoreCase: false,
				invert:     true,
				fixed:      false,
				lineNum:    false,
				pattern:    "Tony",
			},
		},
		expect: []string{
			"Simon Strange",
			"Pete Brown",
			"Stefan Heinz",
			"John Strange",
			"Fred Bloggs",
		},
		expectErr: false,
	}
	t.Run("filter with v=true flag", func(t *testing.T) { test(t, tc) })
}

func Test_filter_with_F_equal_true_flag(t *testing.T) {
	tc := testCase{
		input: filterParams{
			lines: []string{
				"Simon Strange",
				"Pete Brown",
				"Tony Bedford",
				"Stefan Heinz",
				"Tony. Soprano",
				"John Strange",
				"Fred Bloggs",
			},
			args: cmdArgs{
				after:      0,
				before:     0,
				count:      false,
				ignoreCase: false,
				invert:     false,
				fixed:      true,
				lineNum:    false,
				pattern:    "Tony.",
			},
		},
		expect: []string{
			"Tony. Soprano",
		},
		expectErr: false,
	}
	t.Run("filter with F=true flag", func(t *testing.T) { test(t, tc) })
}

func Test_filter_with_n_equal_true_flag(t *testing.T) {
	tc := testCase{
		input: filterParams{
			lines: []string{
				"Simon Strange",
				"Pete Brown",
				"Tony Bedford",
				"Stefan Heinz",
				"Tony Soprano",
				"John Strange",
				"Fred Bloggs",
			},
			args: cmdArgs{
				after:      0,
				before:     0,
				count:      false,
				ignoreCase: false,
				invert:     false,
				fixed:      false,
				lineNum:    true,
				pattern:    "Tony",
			},
		},
		expect: []string{
			"3:Tony Bedford",
			"5:Tony Soprano",
		},
		expectErr: false,
	}
	t.Run("filter with n=true flag", func(t *testing.T) { test(t, tc) })
}

func Test_filter_with_n_equal_true_C_equal_1_flag(t *testing.T) {
	tc := testCase{
		input: filterParams{
			lines: []string{
				"Simon Strange",
				"Pete Brown",
				"Tony Bedford",
				"Stefan Heinz",
				"Tony Soprano",
				"John Strange",
				"Fred Bloggs",
			},
			args: cmdArgs{
				after:      0,
				before:     0,
				count:      false,
				ignoreCase: false,
				invert:     false,
				fixed:      false,
				lineNum:    true,
				pattern:    "Pete",
			},
		},
		expect: []string{
			"1-Simon Strange",
			"2:Pete Brown",
			"3-Tony Bedford",
		},
		expectErr: false,
	}
	t.Run("filter with n=true and C=1 flag", func(t *testing.T) { test(t, tc) })
}

// func Test_filter_with_n_equal_true_C_equal_1_flag(t *testing.T) {
// 	tc := testCase{
// 		input: filterParams{
// 			lines: []string{
// 				"Simon Strange",
// 				"Pete Brown",
// 				"Tony Bedford",
// 				"Stefan Heinz",
// 				"Tony Soprano",
// 				"John Strange",
// 				"Fred Bloggs",
// 			},
// 			args: cmdArgs{
// 				A:       0,
// 				B:       0,
// 				c:       false,
// 				i:       false,
// 				v:       false,
// 				F:       false,
// 				n:       true,
// 				pattern: "Pete",
// 			},
// 		},
// 		expect: []string{
// 			"1-Simon Strange",
// 			"2:Pete Brown",
// 			"3-Tony Bedford",
// 		},
// 		expectErr: false,
// 	}
// 	t.Run("filter with n=true and C=1 flag", func(t *testing.T) { test(t, tc) })
// }
