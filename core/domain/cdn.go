package domain

import (
	"net/url"
	"os"
)

type CDNPath struct {
	Prefix string
	Bucket string
}

func (c CDNPath) URL() url.URL {
	_, isOffline := os.LookupEnv("IS_OFFLINE")

	var cdnServer = ""
	var scheme = ""
	var path = ""

	if isOffline {
		cdnServer = "localhost:4569"
		scheme = "http"
		path = c.Bucket + "/" + c.Prefix
	} else {
		cdnServer = os.Getenv("CDN_SERVER")
		scheme = "https"
		path = c.Prefix
	}

	return url.URL{
		Scheme: scheme,
		Host:   cdnServer,
		Path:   path,
	}
}
