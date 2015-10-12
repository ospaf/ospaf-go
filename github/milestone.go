/*
https://developer.github.com/v3/issues
*/
package githubLib

import (
	"encoding/json"
)

type Milestone struct {
	Url           string
	Html_url      string
	Labels_url    string
	ID            int
	Number        int
	State         string
	Title         string
	Description   string
	Creator       User
	Open_issues   int
	Closed_issues int
	Created_at    string
	Updated_at    string
	Closed_at     string
	Due_on        string
}

func (milestone *Milestone) Marshal() string {
	content, _ := json.MarshalIndent(milestone, "", "  ")
	return string(content)
}

func MilestoneFrom(value string) (milestone Milestone, valid bool) {
	err := json.Unmarshal([]byte(value), &milestone)
	if err != nil {
		valid = false
	} else {
		valid = true
	}

	return milestone, valid
}
