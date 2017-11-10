package exporter

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
)

// gatherData - Collects the data from the API and stores into struct
func (e *Exporter) gatherData() ([]*Datum, *RateLimits, error) {

	data := []*Datum{}

	responses, err := asyncHTTPGets(e.TargetURLs, e.APIToken)

	if err != nil {
		return data, nil, err
	}

	for _, response := range responses {

		// Github can at times present an array, or an object for the same data set.
		// This code checks handles this variation.
		if isArray(response.body) {
			ds := []*Datum{}
			json.Unmarshal(response.body, &ds)
			data = append(data, ds...)
		} else {
			d := new(Datum)
			json.Unmarshal(response.body, &d)
			data = append(data, d)
		}
	}

	// Get stats about branches
	err = getBranchData(e, data)
	if err != nil {
		log.Errorf("Unable to obtain branch data from API, Error: %s", err)
	}

	// Check the API rate data and store as a metric
	rates, err := getRates(e.APIURL, e.APIToken)
	if err != nil {
		log.Errorf("Unable to obtain rate limit data from API, Error: %s", err)
	}

	//return data, rates, err
	return data, rates, nil

}

// isArray simply looks for key details that determine if the JSON response is an array or not.
func isArray(body []byte) bool {

	isArray := false

	for _, c := range body {
		if c == ' ' || c == '\t' || c == '\r' || c == '\n' {
			continue
		}
		isArray = c == '['
		break
	}

	return isArray
}
