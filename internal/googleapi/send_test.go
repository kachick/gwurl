//go:build apitest
// +build apitest

package googleapi

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestPostGoogleAPI(t *testing.T) {
	testCases := []struct {
		description string
		input       App
		want        []string
		ok          bool
	}{
		{
			description: "Google Japanese IME",
			input: App{
				Appid: "{DDCCD2A9-025E-4142-BCEB-F467B88CF830}",
				Ap:    "external-stable-universal",
			},
			want: []string{
				"http://edgedl.me.gvt1.com/edgedl/release2/kjspmop3m4hu2sbbaotsynsgja_2.29.5370.0/GoogleJapaneseInput64-2.29.5370.0.msi",
				"https://edgedl.me.gvt1.com/edgedl/release2/kjspmop3m4hu2sbbaotsynsgja_2.29.5370.0/GoogleJapaneseInput64-2.29.5370.0.msi",
				"http://dl.google.com/release2/kjspmop3m4hu2sbbaotsynsgja_2.29.5370.0/GoogleJapaneseInput64-2.29.5370.0.msi",
				"https://dl.google.com/release2/kjspmop3m4hu2sbbaotsynsgja_2.29.5370.0/GoogleJapaneseInput64-2.29.5370.0.msi",
				"http://www.google.com/dl/release2/kjspmop3m4hu2sbbaotsynsgja_2.29.5370.0/GoogleJapaneseInput64-2.29.5370.0.msi",
				"https://www.google.com/dl/release2/kjspmop3m4hu2sbbaotsynsgja_2.29.5370.0/GoogleJapaneseInput64-2.29.5370.0.msi",
			},
			ok: true,
		},
		{
			description: "Google Chrome",
			input: App{
				Appid: "{8A69D345-D564-463C-AFF1-A69D9E530F96}",
				Ap:    "x64-stable-statsdef_1",
			},
			want: []string{
				"http://edgedl.me.gvt1.com/edgedl/release2/chrome/jowtbwbgcn5nssuzd75wgix2hu_124.0.6367.79/124.0.6367.79_chrome_installer.exe",
				"https://edgedl.me.gvt1.com/edgedl/release2/chrome/jowtbwbgcn5nssuzd75wgix2hu_124.0.6367.79/124.0.6367.79_chrome_installer.exe",
				"http://dl.google.com/release2/chrome/jowtbwbgcn5nssuzd75wgix2hu_124.0.6367.79/124.0.6367.79_chrome_installer.exe",
				"https://dl.google.com/release2/chrome/jowtbwbgcn5nssuzd75wgix2hu_124.0.6367.79/124.0.6367.79_chrome_installer.exe",
				"http://www.google.com/dl/release2/chrome/jowtbwbgcn5nssuzd75wgix2hu_124.0.6367.79/124.0.6367.79_chrome_installer.exe",
				"https://www.google.com/dl/release2/chrome/jowtbwbgcn5nssuzd75wgix2hu_124.0.6367.79/124.0.6367.79_chrome_installer.exe",
			},
			ok: true,
		},
		{
			description: "Unknown Prams",
			input: App{
				Appid: "foo",
				Ap:    "bar",
			},
			want: nil,
			ok:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
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

			dictComp := func(a string, b string) bool {
				return strings.Compare(a, b) == -1
			}

			if diff := cmp.Diff(tc.want, urls, cmpopts.SortSlices(dictComp)); diff != "" {
				t.Errorf("wrong result: %s", diff)
			}
		})
	}
}
