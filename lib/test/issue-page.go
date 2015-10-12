package main

import (
	ospaf "../"
	github "../../github"
	"encoding/json"
	"fmt"
)

func main() {
	pool, err := ospaf.InitPool()
	if err != nil {
		fmt.Println(err)
		return
	}

	owner := "opencontainers"
	repo := "specs"
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/issues/20/comments", owner, repo)
	for page := 1; page != -1; {
		value, code, nextPage, _ := pool.ReadPage(url, page)
		if code != 200 {
			break
		}
		page = nextPage
		var comments []github.Comment
		json.Unmarshal([]byte(value), &comments)
		fmt.Println("There are ", len(comments), "comments")
	}
}
