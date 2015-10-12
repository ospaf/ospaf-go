/*
https://developer.github.com/v3/issues/comments
*/
package githubLib

import (
	"encoding/json"
)

type Comment struct {
	ID         int
	Url        string
	Html_url   string
	Body       string
	User       User
	Created_at string
	Updated_at string
}

func (comment *Comment) Marshal() string {
	content, _ := json.MarshalIndent(comment, "", "  ")
	return string(content)
}

func CommentFrom(value string) (comment Comment, valid bool) {
	err := json.Unmarshal([]byte(value), &comment)
	if err != nil {
		valid = false
	} else {
		valid = true
	}

	return comment, valid
}
