package exporter

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"
)

// getRates obtains the rate limit data for requests against the github API.
// Especially useful when operating without oauth and the subsequent lower cap.
func getRates(baseURL string, token string) (*RateLimits, error) {

	rateEndPoint := ("/rate_limit")
	url := fmt.Sprintf("%s%s", baseURL, rateEndPoint)

	resp, err := getHTTPResponse(url, token)

	defer resp.Body.Close()

	if err != nil {
		return &RateLimits{}, err
	}

	// Triggers if rate-limiting isn't enabled on private Github Enterprise installations
	if resp.StatusCode == 404 {
		return &RateLimits{}, fmt.Errorf("Rate Limiting not enabled in GitHub API")
	}

	limit, err := strconv.ParseFloat(resp.Header.Get("X-RateLimit-Limit"), 64)

	if err != nil {
		return &RateLimits{}, err
	}

	rem, err := strconv.ParseFloat(resp.Header.Get("X-RateLimit-Remaining"), 64)

	if err != nil {
		return &RateLimits{}, err
	}

	reset, err := strconv.ParseFloat(resp.Header.Get("X-RateLimit-Reset"), 64)

	if err != nil {
		return &RateLimits{}, err
	}

	log.Warnf("%d API calls left until %s", int64(rem), time.Unix(int64(reset), 0))

	return &RateLimits{
		Limit:     limit,
		Remaining: rem,
		Reset:     reset,
	}, err

}
