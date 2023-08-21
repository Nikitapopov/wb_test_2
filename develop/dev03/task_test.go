package main

import (
	"reflect"
	"testing"
)

type sortLinesParams struct {
	lines []string
	args  cmdArgs
}

type testCase struct {
	input  sortLinesParams
	expect []string
}

func Test_sortLines_default(t *testing.T) {
	tc := testCase{
		input: sortLinesParams{
			lines: []string{
				"Simon Strange 62",
				"Pete Brown 37",
				"Mark Brown 46",
				"Stefan Heinz 52",
				"Tony Bedford 50",
				"John Strange 51",
				"Fred Bloggs 22",
				"James Bedford 21",
				"Emily Bedford 18",
				"Ana Villamor 44",
				"Alice Villamor 50",
				"Francis Chepstow 56",
			},
			args: cmdArgs{
				k: 1,
				n: false,
				r: false,
				u: false,
			},
		},
		expect: []string{
			"Alice Villamor 50",
			"Ana Villamor 44",
			"Emily Bedford 18",
			"Francis Chepstow 56",
			"Fred Bloggs 22",
			"James Bedford 21",
			"John Strange 51",
			"Mark Brown 46",
			"Pete Brown 37",
			"Simon Strange 62",
			"Stefan Heinz 52",
			"Tony Bedford 50",
		},
	}
	t.Run("default sorting lines", func(t *testing.T) {
		got := sortLines(tc.input.lines, tc.input.args)
		if !reflect.DeepEqual(got, tc.expect) {
			t.Errorf("Expected: %v, got: %v", tc.expect, got)
		}
	})
}

func Test_sortLines_by_second_column(t *testing.T) {
	tc := testCase{
		input: sortLinesParams{
			lines: []string{
				"Simon Strange 62",
				"Pete Brown 37",
				"Mark Brown 46",
				"Stefan Heinz 52",
				"Tony Bedford 50",
				"John Strange 51",
				"Fred Bloggs 22",
				"James Bedford 21",
				"Emily Bedford 18",
				"Ana Villamor 44",
				"Alice Villamor 50",
				"Francis Chepstow 56",
			},
			args: cmdArgs{
				k: 2,
				n: false,
				r: false,
				u: false,
			},
		},
		expect: []string{
			"Emily Bedford 18",
			"James Bedford 21",
			"Tony Bedford 50",
			"Fred Bloggs 22",
			"Pete Brown 37",
			"Mark Brown 46",
			"Francis Chepstow 56",
			"Stefan Heinz 52",
			"John Strange 51",
			"Simon Strange 62",
			"Ana Villamor 44",
			"Alice Villamor 50",
		},
	}
	t.Run("sorting lines by second column", func(t *testing.T) {
		got := sortLines(tc.input.lines, tc.input.args)
		if !reflect.DeepEqual(got, tc.expect) {
			t.Errorf("Expected: %v, got: %v", tc.expect, got)
		}
	})
}

func Test_sortLines_by_number_value(t *testing.T) {
	tc := testCase{
		input: sortLinesParams{
			lines: []string{
				"Simon Strange 62",
				"Pete Brown 37",
				"Mark Brown 46",
				"Stefan Heinz 52",
				"Tony Bedford 50",
				"John Strange 51",
				"Fred Bloggs 22",
				"James Bedford 21",
				"Emily Bedford 18",
				"Ana Villamor 44",
				"Alice Villamor 50",
				"Francis Chepstow 56",
			},
			args: cmdArgs{
				k: 3,
				n: true,
				r: false,
				u: false,
			},
		},
		expect: []string{
			"Emily Bedford 18",
			"James Bedford 21",
			"Fred Bloggs 22",
			"Pete Brown 37",
			"Ana Villamor 44",
			"Mark Brown 46",
			"Tony Bedford 50",
			"Alice Villamor 50",
			"John Strange 51",
			"Stefan Heinz 52",
			"Francis Chepstow 56",
			"Simon Strange 62",
		},
	}
	t.Run("sorting lines by second column", func(t *testing.T) {
		got := sortLines(tc.input.lines, tc.input.args)
		if !reflect.DeepEqual(got, tc.expect) {
			t.Errorf("Expected: %v, got: %v", tc.expect, got)
		}
	})
}

func Test_sortLines_in_reverse_order(t *testing.T) {
	tc := testCase{
		input: sortLinesParams{
			lines: []string{
				"Simon Strange 62",
				"Pete Brown 37",
				"Mark Brown 46",
				"Stefan Heinz 52",
				"Tony Bedford 50",
				"John Strange 51",
				"Fred Bloggs 22",
				"James Bedford 21",
				"Emily Bedford 18",
				"Ana Villamor 44",
				"Alice Villamor 50",
				"Francis Chepstow 56",
			},
			args: cmdArgs{
				k: 1,
				n: false,
				r: true,
				u: false,
			},
		},
		expect: []string{
			"Tony Bedford 50",
			"Stefan Heinz 52",
			"Simon Strange 62",
			"Pete Brown 37",
			"Mark Brown 46",
			"John Strange 51",
			"James Bedford 21",
			"Fred Bloggs 22",
			"Francis Chepstow 56",
			"Emily Bedford 18",
			"Ana Villamor 44",
			"Alice Villamor 50",
		},
	}
	t.Run("sorting lines by second column", func(t *testing.T) {
		got := sortLines(tc.input.lines, tc.input.args)
		if !reflect.DeepEqual(got, tc.expect) {
			t.Errorf("Expected: %v, got: %v", tc.expect, got)
		}
	})
}

func Test_sortLines_unique(t *testing.T) {
	tc := testCase{
		input: sortLinesParams{
			lines: []string{
				"Simon Strange 62",
				"Pete Brown 37",
				"Simon Strange 62",
				"Mark Brown 46",
				"Stefan Heinz 52",
				"Tony Bedford 50",
				"John Strange 51",
				"Fred Bloggs 22",
				"James Bedford 21",
				"Emily Bedford 18",
				"Ana Villamor 44",
				"Alice Villamor 50",
				"Francis Chepstow 56",
			},
			args: cmdArgs{
				k: 1,
				n: false,
				r: false,
				u: true,
			},
		},
		expect: []string{
			"Alice Villamor 50",
			"Ana Villamor 44",
			"Emily Bedford 18",
			"Francis Chepstow 56",
			"Fred Bloggs 22",
			"James Bedford 21",
			"John Strange 51",
			"Mark Brown 46",
			"Pete Brown 37",
			"Simon Strange 62",
			"Stefan Heinz 52",
			"Tony Bedford 50",
		},
	}
	t.Run("sorting lines by second column", func(t *testing.T) {
		got := sortLines(tc.input.lines, tc.input.args)
		if !reflect.DeepEqual(got, tc.expect) {
			t.Errorf("Expected: %v, got: %v", tc.expect, got)
		}
	})
}
