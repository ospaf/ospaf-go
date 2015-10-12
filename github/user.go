/*
https://developer.github.com/v3/users
*/
package githubLib

import (
	"encoding/json"
)

type User struct {
	Login               string
	ID                  int
	Avatar_url          string
	Gravatar_id         string
	Url                 string
	Html_url            string
	Followers_url       string
	Following_url       string
	Gists_url           string
	Starred_url         string
	Subscriptions_url   string
	Organizations_url   string
	Repos_url           string
	Events_url          string
	Received_events_url string
	Type                string
	Site_Admin          bool
}

func (user *User) Marshal() string {
	content, _ := json.MarshalIndent(user, "", "  ")
	return string(content)
}

func UserFrom(value string) (user User, valid bool) {
	err := json.Unmarshal([]byte(value), &user)
	if err != nil {
		valid = false
	} else {
		valid = true
	}

	return user, valid
}
