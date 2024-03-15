package taggedurl

import (
	"net/url"
	"path"
	"slices"
	"strings"

	"golang.org/x/xerrors"
)

type TaggedURL struct {
	Appguid    string
	Ap         string
	Appname    string
	Needsadmin bool
	Filename   string
}

func ParseTaggedURL(likeTaggedUrl string) (TaggedURL, error) {
	// [scheme:][//[userinfo@]host][/]path[?query][#fragment]
	taggedUrl, err := url.ParseRequestURI(likeTaggedUrl)
	if err != nil {
		return TaggedURL{}, xerrors.Errorf("Cannot parse given URL: %w", err)
	}
	if !strings.HasSuffix(taggedUrl.Host, "google.com") {
		return TaggedURL{}, xerrors.Errorf("Given URL looks not a goole: %s", taggedUrl.Host)
	}

	prefixWithQuery, filename := path.Split(taggedUrl.Path)
	// Intentioanlly avoiding path.Split for the getting nth element. Not the prefix and last
	dirs := strings.Split(prefixWithQuery, "/")
	qsi := slices.IndexFunc(dirs, func(dir string) bool { return strings.Contains(dir, "appguid") })
	qs := dirs[qsi]
	query, err := url.ParseQuery(qs)
	if err != nil {
		return TaggedURL{}, xerrors.Errorf("Cannot parse given query: %w", err)
	}
	appguid, ok := query["appguid"]
	if !ok {
		return TaggedURL{}, xerrors.Errorf("No appguid: %s", appguid)
	}
	ap, ok := query["ap"]
	if !ok {
		return TaggedURL{}, xerrors.Errorf("No ap: %s", ap)
	}
	appname, ok := query["appname"]
	if !ok {
		return TaggedURL{}, xerrors.Errorf("No appname: %s", appname)
	}
	needsadmin, ok := query["needsadmin"]
	if !ok {
		return TaggedURL{}, xerrors.Errorf("No needsadmin: %s", needsadmin)
	}

	return TaggedURL{
		Appguid:    appguid[0],
		Ap:         ap[0],
		Appname:    appname[0],
		Needsadmin: needsadmin[0] == "true",
		Filename:   filename,
	}, nil
}
