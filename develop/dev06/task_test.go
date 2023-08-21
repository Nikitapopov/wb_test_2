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
	input  filterParams
	expect []string
}

func test(t *testing.T, tc testCase) {
	got := cut(tc.input.lines, tc.input.args)
	if !reflect.DeepEqual(got, tc.expect) {
		t.Errorf("Expected: %v, got: %v", tc.expect, got)
	}
}

func Test_cut_by_default_options(t *testing.T) {
	tc := testCase{
		input: filterParams{
			lines: []string{
				"Stefan Heinz Marcus Mirele",
				"Tony	Soprano",
				"John,Strange	123",
				"Brr",
				"",
				"Fred,Bloggs,Frog,Eooo",
			},
			args: cmdArgs{
				fields:    []int{1},
				delimiter: "\t",
				separated: false,
			},
		},
		expect: []string{
			"Stefan Heinz Marcus Mirele",
			"Tony",
			"John,Strange",
			"Brr",
			"",
			"Fred,Bloggs,Frog,Eooo",
		},
	}
	t.Run("cut by default options", func(t *testing.T) { test(t, tc) })
}

func Test_cut_by_default_options_from_empty_list(t *testing.T) {
	tc := testCase{
		input: filterParams{
			lines: []string{},
			args: cmdArgs{
				fields:    []int{1},
				delimiter: "\t",
				separated: false,
			},
		},
		expect: []string{},
	}
	t.Run("cut by default options from empty list", func(t *testing.T) { test(t, tc) })
}

func Test_cut_by_fields_1_2(t *testing.T) {
	tc := testCase{
		input: filterParams{
			lines: []string{
				"Stefan Heinz Marcus Mirele",
				"Tony	Soprano",
				"John,Strange	123",
				"Brr",
				"",
				"Fred,Bloggs,Frog,Eooo",
			},
			args: cmdArgs{
				fields:    []int{1, 2},
				delimiter: "\t",
				separated: false,
			},
		},
		expect: []string{
			"Stefan Heinz Marcus Mirele",
			"Tony	Soprano",
			"John,Strange	123",
			"Brr",
			"",
			"Fred,Bloggs,Frog,Eooo",
		},
	}
	t.Run("cut by fields=[1,2]", func(t *testing.T) { test(t, tc) })
}

func Test_cut_by_fields_2(t *testing.T) {
	tc := testCase{
		input: filterParams{
			lines: []string{
				"Stefan Heinz Marcus Mirele",
				"Tony	Soprano",
				"John,Strange	123",
				"Brr",
				"",
				"Fred,Bloggs,Frog,Eooo",
			},
			args: cmdArgs{
				fields:    []int{2},
				delimiter: "\t",
				separated: false,
			},
		},
		expect: []string{
			"",
			"Soprano",
			"123",
			"",
			"",
			"",
		},
	}
	t.Run("cut by fields=[2]", func(t *testing.T) { test(t, tc) })
}

func Test_cut_by_space_delimiter_fields_2(t *testing.T) {
	tc := testCase{
		input: filterParams{
			lines: []string{
				"Stefan Heinz Marcus Mirele",
				"Tony	Soprano",
				"John,Strange	123",
				"Brr",
				"",
				"Fred,Bloggs,Frog,Eooo",
			},
			args: cmdArgs{
				fields:    []int{2},
				delimiter: " ",
				separated: false,
			},
		},
		expect: []string{
			"Heinz",
			"",
			"",
			"",
			"",
			"",
		},
	}
	t.Run(`cut by delimiter=" " fields=[2]`, func(t *testing.T) { test(t, tc) })
}

func Test_cut_by_comma_delimiter(t *testing.T) {
	tc := testCase{
		input: filterParams{
			lines: []string{
				"Stefan Heinz Marcus Mirele",
				"Tony	Soprano",
				"John,Strange	123",
				"Brr",
				"",
				"Fred,Bloggs,Frog,Eooo",
			},
			args: cmdArgs{
				fields:    []int{1},
				delimiter: ",",
				separated: false,
			},
		},
		expect: []string{
			"Stefan Heinz Marcus Mirele",
			"Tony	Soprano",
			"John",
			"Brr",
			"",
			"Fred",
		},
	}
	t.Run(`cut by delimiter=","`, func(t *testing.T) { test(t, tc) })
}

func Test_cut_by_separates_true_comma_delimiter_fields_1_2(t *testing.T) {
	tc := testCase{
		input: filterParams{
			lines: []string{
				"Stefan Heinz Marcus Mirele",
				"Tony	Soprano",
				"John,Strange	123",
				"Brr",
				"",
				"Fred,Bloggs,Frog,Eooo",
			},
			args: cmdArgs{
				fields:    []int{1, 2},
				delimiter: ",",
				separated: true,
			},
		},
		expect: []string{
			"John,Strange	123",
			"Fred,Bloggs",
		},
	}
	t.Run(`cut by separeted=true, delimiter="," and fields=[1,2]`, func(t *testing.T) { test(t, tc) })
}
