package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
)

var (
	// Used in goreleaser
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

	original := os.Args[1]
	permalink, err := url.ParseRequestURI(original)

	if err != nil {
		log.Fatalf("Cannot parse given URL: %+v", err)
	}

	if !strings.HasSuffix(permalink.Host, "google.com") {
		log.Fatalf("Given URL looks not a goole: %s", permalink.Host)
	}
	query := permalink.Query()
	appguid, ok := query["appguid"]
	if !ok {
		log.Fatalf("No appguid: %s", appguid)
	}
	ap, ok := query["ap"]
	if !ok {
		log.Fatalf("No ap: %s", ap)
	}

	fmt.Printf("%s, %s, %+v\n", appguid, ap, query)
}
