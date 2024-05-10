//go:build apitest
// +build apitest

package googleapi

import (
	"regexp"
	"slices"
	"testing"
)

func TestPostGoogleAPI(t *testing.T) {
	testCases := []struct {
		description string
		input       App
		want        string
		ok          bool
	}{
		{
			description: "Google Japanese IME",
			input: App{
				Appid: "{DDCCD2A9-025E-4142-BCEB-F467B88CF830}",
				Ap:    "external-stable-universal",
			},
			want: "https?://[a-zA-Z0-9/.-]+_[0-9.]+/GoogleJapaneseInput[a-zA-Z0-9/.-]+\\.(exe|msi)",
			ok:   true,
		},
		{
			description: "Google Chrome",
			input: App{
				Appid: "{8A69D345-D564-463C-AFF1-A69D9E530F96}",
				Ap:    "x64-stable-statsdef_1",
			},
			want: "https?://[a-zA-Z0-9/\\.-]+_[0-9\\.]+/[a-zA-Z0-9/\\.-]+_chrome_installer\\.(exe|msi)",
			ok:   true,
		},
		{
			description: "Unknown parameters",
			input: App{
				Appid: "foo",
				Ap:    "bar",
			},
			ok: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			var wantPattern *regexp.Regexp
			var err error
			if tc.ok {
				wantPattern, err = regexp.Compile(tc.want)
				if err != nil {
					t.Fatalf("unexpected error happened: %v", err)
				}
			}

			resp, err := PostGoogleAPI(Os{
				Platform:     "win",
				Version:      "10",
				Architecture: "x64",
			}, tc.input)
			if err != nil {
				if tc.ok {
					t.Fatalf("unexpected error happened: %v", err)
				} else {
					return
				}
			}

			urls, err := GetPermalinks(resp)
			if err != nil {
				if tc.ok {
					t.Fatalf("unexpected error happened: %v", err)
				} else {
					return
				}
			}

			if !tc.ok {
				t.Fatalf("expected error did not happen")
			}
			if slices.ContainsFunc(urls, func(u string) bool { return !wantPattern.MatchString(u) }) {
				t.Errorf("returned urls contain unexpected pattern: %v", urls)
			}
		})
	}
}
