/*
https://developer.github.com/v3/rate_limit/
*/
package githubLib

import (
	"encoding/json"
)

type RateLimitUnit struct {
	Limit     int
	Remaining int
	Reset     int64
}

type RateLimitResource struct {
	Core   RateLimitUnit
	Search RateLimitUnit
}

type RateLimit struct {
	Resources RateLimitResource
	Rate      RateLimitUnit
}

func (rateLimit *RateLimit) Marshal() string {
	content, _ := json.MarshalIndent(rateLimit, "", "  ")
	return string(content)
}

func RateLimitFrom(value string) (rateLimit RateLimit, valid bool) {
	err := json.Unmarshal([]byte(value), &rateLimit)
	if err != nil {
		valid = false
	} else {
		valid = true
	}

	return rateLimit, valid
}
