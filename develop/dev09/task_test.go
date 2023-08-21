package main

import (
	"os"
	"testing"
)

type filterParams struct {
	args cmdArgs
}

type testCase struct {
	input     filterParams
	expectErr bool
}

func test(t *testing.T, tc testCase) {
	err := getWebsite(tc.input.args)
	if err != nil && !tc.expectErr {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	if err != nil && !tc.expectErr {
		t.Errorf("Expected error: %v", tc.expectErr)
		return
	}

	filepath := "data/go.dev/play.html1"
	_, err = os.Stat(filepath)
	if err != nil {
		t.Error("website's file not found")
	}
}

func Test_wget(t *testing.T) {
	tc := testCase{
		input: filterParams{
			args: cmdArgs{
				uri:       "https://go.dev/play/",
				recursive: true,
				depth:     0,
			},
		},
	}
	t.Run("get website without recursive download", func(t *testing.T) { test(t, tc) })
}
