package googleapi

import (
	"bytes"
	"encoding/xml"
	"io"
	"net/http"
	"net/url"
	"slices"

	"golang.org/x/xerrors"
)

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

type Os struct {
	Platform     string `xml:"platform,attr"`
	Version      string `xml:"version,attr"`
	Architecture string `xml:"arch,attr"`
}

type App struct {
	Appid       string   `xml:"appid,attr"`
	Version     string   `xml:"version,attr"`
	Ap          string   `xml:"ap,attr"`
	UpdateCheck struct{} `xml:"updatecheck"`
}

type Request struct {
	XMLName  xml.Name `xml:"request"`
	Protocol string   `xml:"protocol,attr"`
	Os       Os       `xml:"os"`
	App      App      `xml:"app"`
}

func BuildGoogleApiPostXml(os Os, app App) ([]byte, error) {
	full := []byte(xml.Header)
	body, err := xml.Marshal(Request{
		Protocol: "3.0",
		Os:       os,
		App:      app,
	})
	full = append(full, body...)

	return full, err
}

func PostGoogleAPI(os Os, app App) (Response, error) {
	body, err := BuildGoogleApiPostXml(os, app)
	if err != nil {
		return Response{}, xerrors.Errorf("Error creating request: %w", err)
	}

	req, err := http.NewRequest("POST", "https://update.googleapis.com/service/update2", bytes.NewBuffer(body))
	if err != nil {
		return Response{}, xerrors.Errorf("Error sending request: %w", err)
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

func GetPermalinks(resp Response) ([]string, error) {
	installerActionIdx := slices.IndexFunc(resp.App.UpdateCheck.Manifest.Actions, func(a Action) bool {
		return a.Event == "install"
	})
	if installerActionIdx < 0 {
		return nil, xerrors.Errorf("api didn't return installer information: %v", installerActionIdx)
	}
	installerFilename := resp.App.UpdateCheck.Manifest.Actions[installerActionIdx].Run

	permalinks := make([]string, 0, len(resp.App.UpdateCheck.Urls))
	for _, u := range resp.App.UpdateCheck.Urls {
		permalink, err := url.JoinPath(u.Codebase, installerFilename)
		if err != nil {
			return nil, xerrors.Errorf("Cannot build final link with the result: %w", err)
		}
		permalinks = append(permalinks, permalink)
	}
	return permalinks, nil
}
