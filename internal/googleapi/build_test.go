package googleapi

import (
	"testing"
)

func TestBuildGoogleApiPostXml(t *testing.T) {
	t.Run("Build XML Request", func(t *testing.T) {
		_, err := BuildGoogleApiPostXml(Os{
			Platform:     "win",
			Version:      "10",
			Architecture: "x64",
		}, App{
			Appid: "{DDCCD2A9-025E-4142-BCEB-F467B88CF830}",
			Ap:    "external-stable-universal",
		})
		if err != nil {
			t.Fatalf("unexpected error happned: %+v", err)
		}
	})
}
