package main

import (
	github "../../github"
	ospaf "../../lib"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Owner string
	Repo  string
}

func loadComments(url string, pool ospaf.Pool, config Config, number int) {
	fmt.Println("Start to load: ", url)

	var commentList []github.Comment
	paras := make(map[string]string)
	for page := 1; page != -1; {
		paras["page"] = strconv.Itoa(page)
		value, code, nextPage, _ := pool.ReadPage(url, paras)
		if code != 200 {
			break
		}
		page = nextPage

		var comments []github.Comment
		json.Unmarshal([]byte(value), &comments)
		for index := 0; index < len(comments); index++ {
			commentList = append(commentList, comments[index])
		}
	}

	fileUrl := fmt.Sprintf("data/comment-of-issue-%s-%s-%d", config.Owner, config.Repo, number)
	content, _ := json.MarshalIndent(commentList, "", "  ")
	fout, err := os.Create(fileUrl)
	if err != nil {
		fmt.Println(fileUrl, err)
	} else {
		fout.WriteString(string(content))
		fmt.Println("Save ", url, " to ", fileUrl)
		fout.Close()
	}
}

func main() {
	pool, err := ospaf.InitPool()
	if err != nil {
		fmt.Println(err)
		return
	}

	ospaf.PreparePath("data", "")

	content, err := ospaf.ReadFile("./config.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	var config Config
	json.Unmarshal([]byte(content), &config)

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/issues", config.Owner, config.Repo)
	paras := make(map[string]string)
	paras["state"] = "all"

	var issueList []github.Issue
	for page := 1; page != -1; {
		paras["page"] = strconv.Itoa(page)
		value, code, nextPage, _ := pool.ReadPage(url, paras)
		if code != 200 {
			break
		}
		page = nextPage

		var issues []github.Issue
		json.Unmarshal([]byte(value), &issues)

		for index := 0; index < len(issues); index++ {
			if issues[index].Pull_request.Url == "" {
				fmt.Println("Pull request, drop")
				continue
			}
			issueList = append(issueList, issues[index])
			commentUrl := issues[index].Comments_url
			loadComments(commentUrl, pool, config, issues[index].Number)
		}
	}

	fileUrl := fmt.Sprintf("data/issue-of-%s-%s", config.Owner, config.Repo)
	ilContent, _ := json.MarshalIndent(issueList, "", "  ")
	fout, err := os.Create(fileUrl)
	if err != nil {
		fmt.Println(fileUrl, err)
	} else {
		fout.WriteString(string(ilContent))
		fmt.Println("Save ", url, " to ", fileUrl)
		fout.Close()
	}

}
