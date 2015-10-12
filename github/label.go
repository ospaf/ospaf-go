/*
https://developer.github.com/v3/issues
*/
package githubLib

import (
	"encoding/json"
)

type Label struct {
	Url   string
	Name  string
	Color string
}

func (label *Label) Marshal() string {
	content, _ := json.MarshalIndent(label, "", "  ")
	return string(content)
}

func LabelFrom(value string) (label Label, valid bool) {
	err := json.Unmarshal([]byte(value), &label)
	if err != nil {
		valid = false
	} else {
		valid = true
	}

	return label, valid
}
