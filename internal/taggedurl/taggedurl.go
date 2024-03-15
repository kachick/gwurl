package taggedurl

import (
	"log"
	"net/url"
	"path"
	"slices"
	"strings"
)

type TaggedURL struct {
	Appguid    string
	Ap         string
	Appname    string
	Needsadmin bool
	Filename   string
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
		Appguid:    appguid[0],
		Ap:         ap[0],
		Appname:    appname[0],
		Needsadmin: needsadmin[0] == "true",
		Filename:   filename,
	}
}
