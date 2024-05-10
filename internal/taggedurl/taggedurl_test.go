package taggedurl

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParseTaggedURL(t *testing.T) {
	testCases := []struct {
		description string
		input       string
		want        TaggedURL
		ok          bool
	}{
		{
			description: "Google Japanese IME",
			input:       "https://dl.google.com/tag/s/appguid%3D%7BDDCCD2A9-025E-4142-BCEB-F467B88CF830%7D%26iid%3D%7BBF6725E7-4A15-87B7-24D5-0ADABEC753EE%7D%26lang%3Dja%26browser%3D4%26usagestats%3D0%26appname%3DGoogle%2520%25E6%2597%25A5%25E6%259C%25AC%25E8%25AA%259E%25E5%2585%25A5%25E5%258A%259B%26needsadmin%3Dtrue%26ap%3Dexternal-stable-universal/japanese-ime/GoogleJapaneseInputSetup.exe",
			want: TaggedURL{
				Appguid:    "{DDCCD2A9-025E-4142-BCEB-F467B88CF830}",
				Ap:         "external-stable-universal",
				Appname:    "Google 日本語入力",
				Needsadmin: true,
				Filename:   "GoogleJapaneseInputSetup.exe",
			},
			ok: true,
		},
		{
			description: "Google Chrome",
			input:       "https://dl.google.com/tag/s/appguid%3D%7B8A69D345-D564-463C-AFF1-A69D9E530F96%7D%26iid%3D%7BF99B66DF-85A3-C043-7205-7917C2AA7AAB%7D%26lang%3Dja%26browser%3D4%26usagestats%3D1%26appname%3DGoogle%2520Chrome%26needsadmin%3Dprefers%26ap%3Dx64-stable-statsdef_1%26installdataindex%3Dempty/update2/installers/ChromeSetup.exe",
			want: TaggedURL{
				Appguid:    "{8A69D345-D564-463C-AFF1-A69D9E530F96}",
				Ap:         "x64-stable-statsdef_1",
				Appname:    "Google Chrome",
				Needsadmin: false,
				Filename:   "ChromeSetup.exe",
			},
			ok: true,
		},
		{
			description: "Unknown URL",
			input:       "https://example.cpm/dir/foobar.exe",
			want:        TaggedURL{},
			ok:          false,
		},
		{
			description: "Not a URL",
			input:       ":)",
			want:        TaggedURL{},
			ok:          false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			parsed, err := ParseTaggedURL(tc.input)

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

			if diff := cmp.Diff(tc.want, parsed); diff != "" {
				t.Errorf("wrong result: %s", diff)
			}
		})
	}
}
