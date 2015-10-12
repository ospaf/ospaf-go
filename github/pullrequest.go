/*
https://developer.github.com/v3/issues
*/
package githubLib

import (
	"encoding/json"
)

type PullRequest struct {
	Url       string
	Html_url  string
	Diff_url  string
	Patch_url string
}

func (pullRequest *PullRequest) Marshal() string {
	content, _ := json.MarshalIndent(pullRequest, "", "  ")
	return string(content)
}

func PullRequestFrom(value string) (pullRequest PullRequest, valid bool) {
	err := json.Unmarshal([]byte(value), &pullRequest)
	if err != nil {
		valid = false
	} else {
		valid = true
	}

	return pullRequest, valid
}
