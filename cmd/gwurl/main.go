package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/kachick/gwurl/internal/googleapi"
	"github.com/kachick/gwurl/internal/taggedurl"
)

var (
	// Might be used in goreleaser
	version = "dev"
	commit  = "none"

	revision = "rev"
)

func main() {
	versionFlag := flag.Bool("version", false, "print the version of this program")

	const usage = `Usage: gwurl [OPTIONS] [URL]

$ gwurl "$Windows_Installer_URL_That_Provided_By_Google"
$ gwurl --version
`

	flag.Usage = func() {
		// https://github.com/golang/go/issues/57059#issuecomment-1336036866
		fmt.Printf("%s", usage+"\n\n")
		fmt.Println("Usage of command:")
		flag.PrintDefaults()
	}

	if len(commit) >= 7 {
		revision = commit[:7]
	}
	version := fmt.Sprintf("%s\n", "gwurl"+" "+version+" "+"("+revision+")")

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	flag.Parse()
	if *versionFlag {
		fmt.Println(version)
		return
	}

	taggedUrl := os.Args[1]
	parsed, err := taggedurl.ParseTaggedURL(taggedUrl)
	if err != nil {
		log.Fatalf("Cannot parse given URL: %+v", err)
	}
	fmt.Printf("Given Tagged URL Inspection: %+v\n", parsed)

	resp, err := googleapi.PostGoogleAPI(googleapi.Os{
		Platform:     "win",
		Version:      "10",
		Architecture: "x64",
	}, googleapi.App{
		Appid: parsed.Appguid,
		Ap:    parsed.Ap,
	})
	if err != nil {
		log.Fatalf("Cannot ask to Google API: %+v", err)
	}

	permalinks, err := googleapi.GetPermalinks(resp)
	if err != nil {
		log.Fatalf("Cannot ask to Google API: %+v", err)
	}

	for _, p := range permalinks {
		fmt.Println(p)
	}
}
