package googleapi

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"

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

type GoogleApiOs struct {
	Platform     string
	Version      string
	Architecture string
}

type GoogleApiApp struct {
	Appid string
	Ap    string
}

func BuildGoogleApiPostXml(apiOs GoogleApiOs, apiApp GoogleApiApp) string {
	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<request protocol="3.0">
  <os platform="%s" version="%s" arch="%s" />
  <app appid="%s" version="" ap="%s">
    <updatecheck />
  </app>
</request>
`, apiOs.Platform, apiOs.Version, apiOs.Architecture, apiApp.Appid, apiApp.Ap)
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
