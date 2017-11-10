package exporter

import (
	"net/http"

	"github.com/infinityworks/github-exporter/config"
	"github.com/prometheus/client_golang/prometheus"
)

// Exporter is used to store Metrics data and embeds the config struct.
// This is done so that the relevant functions have easy access to the
// user defined runtime configuration when the Collect method is called.
type Exporter struct {
	APIMetrics map[string]*prometheus.Desc
	config.Config
}

// Data is used to store an array of Datums.
// This is useful for the JSON array detection
type Data []Datum

// Datum is used to store data from all the relevant endpoints in the API
type Datum struct {
	Name  string `json:"name"`
	Owner struct {
		Login string `json:"login"`
	} `json:"owner"`
	Private    bool    `json:"private"`
	Forks      float64 `json:"forks"`
	Stars      float64 `json:"stargazers_count"`
	OpenIssues float64 `json:"open_issues"`
	Watchers   float64 `json:"subscribers_count"`
	Branches   float64
	Size       float64 `json:"size"`
}

type Branch struct {
	Name          string `json:"name"`
	Protected     bool   `json:"protected"`
	ProtectionUrl string `json:"protection_url"`
}

// RateLimits is used to store rate limit data into a struct
// This data is later represented as a metric, captured at the end of a scrape
type RateLimits struct {
	Limit     float64
	Remaining float64
	Reset     float64
}

// Response struct is used to store http.Response and associated data
type Response struct {
	url      string
	response *http.Response
	body     []byte
	err      error
}
