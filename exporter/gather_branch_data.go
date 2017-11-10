package exporter

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

func getBranchData(e *Exporter, data []*Datum) error {

	branchURLs := generateBranchUrls(data, e)

	responses, err := asyncHTTPGets(branchURLs, e.APIToken)
	if err != nil {
		return err
	}

	for _, response := range responses {
		repository := extractRepositoryName(response.url)
		numberOfBranches, err := countNumberOfBranches(response.body, response.response.Header, e.APIToken)
		if err != nil {
			return err
		}

		for _, repo := range data {
			if repo.Name == repository {
				repo.Branches += float64(numberOfBranches)
			}
		}
	}

	return nil
}

func generateBranchUrls(data []*Datum, e *Exporter) []string {
	branchURLs := []string{}
	for _, repo := range data {
		url := fmt.Sprintf("%s/repos/%s/%s/branches", e.APIURL, repo.Owner.Login, repo.Name)
		branchURLs = append(branchURLs, url)
	}
	return branchURLs
}

func extractRepositoryName(url string) string {
	r := regexp.MustCompile("^.*?/repos/.*?/(.*?)/branches$")
	if r.MatchString(url) {
		matches := r.FindStringSubmatch(url)
		return matches[1]
	}
	return ""
}

func countNumberOfBranches(body []byte, headers http.Header, token string) (int, error) {
	numberOfBranches := 0

	branches := []*Branch{}
	json.Unmarshal(body, &branches)

	numberOfBranches += len(branches)

	linkHeader := getLinkHeader(headers)

	if hasNextLink(linkHeader) {
		linkUrl := extractNextLinkUrl(linkHeader)

		response, err := getHTTPResponse(linkUrl, token)
		if err != nil {
			return 0, err
		}

		// Read the body to a byte array so it can be used elsewhere
		body, err := ioutil.ReadAll(response.Body)

		count, err := countNumberOfBranches(body, response.Header, token)
		if err != nil {
			return 0, err
		}

		numberOfBranches += count

		defer response.Body.Close()

		if err != nil {
			return 0, err
		}

		// Triggers if a user specifies an invalid or not visible repository
		if response.StatusCode == 404 {
			return 0, fmt.Errorf("Error: Received 404 status from Github API, ensure the repsository URL is correct. If it's a privare repository, also check the oauth token is correct")
		}
	}

	return numberOfBranches, nil
}

func getLinkHeader(headers http.Header) string {
	return headers.Get("Link")
}

func hasNextLink(linkHeader string) bool {
	return strings.Contains(linkHeader, "rel=\"next\"")
}

func extractNextLinkUrl(header string) string {
	r := regexp.MustCompile("<(.*?)>; rel=\"next\"")
	matches := r.FindStringSubmatch(header)
	return matches[1]
}
