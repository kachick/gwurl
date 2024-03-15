package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"path"
	"slices"
	"strings"
)

var (
	// Used in goreleaser
	version = "dev"
	commit  = "none"

	revision = "rev"
)

type TaggedURL struct {
	appguid    string
	ap         string
	appname    string
	needsadmin bool
	filename   string
}

func ParseTaggedURL(taggedUrl string) TaggedURL {
	// [scheme:][//[userinfo@]host][/]path[?query][#fragment]
	permalink, err := url.ParseRequestURI(taggedUrl)
	if err != nil {
		log.Fatalf("Cannot parse given URL: %+v", err)
	}
	if !strings.HasSuffix(permalink.Host, "google.com") {
		log.Fatalf("Given URL looks not a goole: %s", permalink.Host)
	}

	prefixWithQuery, filename := path.Split(permalink.Path)
	// Intentioanlly avoiding path.Split for the getting nth element. Not the prefix and last
	dirs := strings.Split(prefixWithQuery, "/")
	qsi := slices.IndexFunc(dirs, func(dir string) bool { return strings.Contains(dir, "appguid") })
	qs := dirs[qsi]
	query, err := url.ParseQuery(qs)
	if err != nil {
		log.Fatalf("Cannot parse given query: %+v", err)
	}
	appguid, ok := query["appguid"]
	if !ok {
		log.Fatalf("No appguid: %s", appguid)
	}
	ap, ok := query["ap"]
	if !ok {
		log.Fatalf("No ap: %s", ap)
	}
	appname, ok := query["appname"]
	if !ok {
		log.Fatalf("No appname: %s", ap)
	}
	needsadmin, ok := query["needsadmin"]
	if !ok {
		log.Fatalf("No needsadmin: %s", needsadmin)
	}

	return TaggedURL{
		appguid:    appguid[0],
		ap:         ap[0],
		appname:    appname[0],
		needsadmin: needsadmin[0] == "true",
		filename:   filename,
	}
}

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
	parsed := ParseTaggedURL(taggedUrl)

	// fmt.Printf("%s, %s, %+v\n", appguid, ap, query)
	fmt.Printf("%+v\n", parsed)
}
