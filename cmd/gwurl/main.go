package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"path"
	"slices"
	"strings"

	"bytes"
	"encoding/xml"
	"net/http"

	"golang.org/x/xerrors"
)

var (
	// Might be used in goreleaser
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

type Action struct {
	Event string `xml:"event,attr"`
	Run   string `xml:"run,attr"`
}

// NOTE: Use `InnerXML string `xml:",innerxml"“ to inspect unknown fields: https://stackoverflow.com/a/38509722
type Response struct {
	XMLName xml.Name `xml:"response"`
	App     struct {
		Status      string `xml:"status"`
		UpdateCheck struct {
			Manifest struct {
				Version string   `xml:"version,attr"`
				Actions []Action `xml:"actions>action"`
			} `xml:"manifest"`

			Urls []struct {
				Codebase string `xml:"codebase,attr"`
			} `xml:"urls>url"`
		} `xml:"updatecheck"`
	} `xml:"app"`
}

type GoogleApiOs struct {
	platform     string
	version      string
	architecture string
}

type GoogleApiApp struct {
	appid string
	ap    string
}

func BuildGoogleApiPostXml(apiOs GoogleApiOs, apiApp GoogleApiApp) string {
	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<request protocol="3.0">
  <os platform="%s" version="%s" arch="%s" />
  <app appid="%s" version="" ap="%s">
    <updatecheck />
  </app>
</request>
`, apiOs.platform, apiOs.version, apiOs.architecture, apiApp.appid, apiApp.ap)
}

func PostGoogleAPI(apiOs GoogleApiOs, apiApp GoogleApiApp) (Response, error) {
	body := []byte(BuildGoogleApiPostXml(apiOs, apiApp))

	req, err := http.NewRequest("POST", "https://update.googleapis.com/service/update2", bytes.NewBuffer(body))
	if err != nil {
		return Response{}, xerrors.Errorf("Error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/xml; charset=utf-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return Response{}, xerrors.Errorf("Error posting request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return Response{}, xerrors.Errorf("Error reading response body: %w", err)
	}

	var responseObject Response
	err = xml.Unmarshal(respBody, &responseObject)
	if err != nil {
		return Response{}, xerrors.Errorf("Error unmarshalling response: %w", err)
	}

	return responseObject, nil
}

func ParseTaggedURL(likeTaggedUrl string) TaggedURL {
	// [scheme:][//[userinfo@]host][/]path[?query][#fragment]
	taggedUrl, err := url.ParseRequestURI(likeTaggedUrl)
	if err != nil {
		log.Fatalf("Cannot parse given URL: %+v", err)
	}
	if !strings.HasSuffix(taggedUrl.Host, "google.com") {
		log.Fatalf("Given URL looks not a goole: %s", taggedUrl.Host)
	}

	prefixWithQuery, filename := path.Split(taggedUrl.Path)
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
	fmt.Printf("%+v\n", parsed)

	resp, err := PostGoogleAPI(GoogleApiOs{
		platform:     "win",
		version:      "10",
		architecture: "x64",
	}, GoogleApiApp{
		appid: parsed.appguid,
		ap:    parsed.ap,
	})
	if err != nil {
		log.Fatalf("Cannot ask to Google API: %+v", err)
	}

	installerActionIdx := slices.IndexFunc(resp.App.UpdateCheck.Manifest.Actions, func(a Action) bool {
		return a.Event == "install"
	})
	installerFilename := resp.App.UpdateCheck.Manifest.Actions[installerActionIdx].Run

	for _, u := range resp.App.UpdateCheck.Urls {
		permalink, err := url.JoinPath(u.Codebase, installerFilename)
		if err != nil {
			log.Fatalf("Cannot build final link with the result: %+v", err)
		}
		fmt.Println(permalink)
	}
}
