/*
https://developer.github.com/v3/issues
*/
package githubLib

import (
	"encoding/json"
)

type Issue struct {
	ID           int
	Url          string
	Labels_url   string
	Comments_url string
	Events_url   string
	html_url     string
	Number       int
	State        string
	Title        string
	Body         string
	User         User
	Labels       []Label
	Assignee     User
	Milestone    Milestone
	Locked       bool
	Comments     int
	Pull_request PullRequest
	Closed_at    string
	Create_at    string
	Updated_at   string
}

func (issue *Issue) Marshal() string {
	content, _ := json.MarshalIndent(issue, "", "  ")
	return string(content)
}

func IssueFrom(value string) (issue Issue, valid bool) {
	err := json.Unmarshal([]byte(value), &issue)
	if err != nil {
		valid = false
	} else {
		valid = true
	}

	return issue, valid
}
