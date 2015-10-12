package main

import (
	github "../../github"
	ospaf "../../lib"
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

type Config struct {
	Owner string
	Repo  string
}

type Hacker struct {
	Login string
	Value float64
}

type HackerList []Hacker

func (list HackerList) Len() int {
	return len(list)
}

func (list HackerList) Less(i, j int) bool {
	if list[i].Value < list[j].Value {
		return true
	} else if list[i].Value > list[j].Value {
		return false
	} else {
		return list[i].Login < list[j].Login
	}
}

func (list HackerList) Swap(i, j int) {
	var temp Hacker = list[i]
	list[i] = list[j]
	list[j] = temp
}

func openIssueReport(issueList []github.Issue, fileUrl string) {
	issueMap := make(map[string]float64)
	for index := 0; index < len(issueList); index++ {
		issue := issueList[index]
		v, ok := issueMap[issue.User.Login]
		if ok {
			v = v + 1.0
		} else {
			v = 1.0
		}
		issueMap[issue.User.Login] = v
	}
	var hackerList HackerList
	wholeIssueValue := 0.0
	wholeReporter := 0
	for key, value := range issueMap {
		var hacker Hacker
		hacker.Login = key
		hacker.Value = value
		hackerList = append(hackerList, hacker)
		wholeIssueValue += hacker.Value
		wholeReporter += 1
	}
	sort.Sort(sort.Reverse(hackerList))

	content := fmt.Sprintf("There are %d issues reported by %d hackers:\n", int(wholeIssueValue), wholeReporter)
	for index := 0; index < wholeReporter; index++ {
		content += fmt.Sprintf("%d. %s reports %d issues, %.2f%%.\n", index+1, hackerList[index].Login, int(hackerList[index].Value), 100*hackerList[index].Value/wholeIssueValue)
	}
	fmt.Println(content)

	fout, err := os.Create(fileUrl)
	if err != nil {
		fmt.Println(fileUrl, err)
	} else {
		fout.WriteString(string(content))
		fmt.Println("Generate ", fileUrl)
		fout.Close()
	}
}

func issueCommentReport(commentList []github.Comment, fileUrl string) {
	commentMap := make(map[string]float64)
	for index := 0; index < len(commentList); index++ {
		comment := commentList[index]
		v, ok := commentMap[comment.User.Login]
		if ok {
			v++
		} else {
			v = 1
		}
		commentMap[comment.User.Login] = v
	}
	var hackerList HackerList
	wholeCommentValue := 0.0
	wholeCommentator := 0
	for key, value := range commentMap {
		var hacker Hacker
		hacker.Login = key
		hacker.Value = value
		hackerList = append(hackerList, hacker)
		wholeCommentValue += hacker.Value
		wholeCommentator += 1
	}
	sort.Sort(sort.Reverse(hackerList))

	content := fmt.Sprintf("There are %d comments reported by %d hackers:\n", int(wholeCommentValue), wholeCommentator)
	for index := 0; index < wholeCommentator; index++ {
		content += fmt.Sprintf("%d. %s comments %d times, %.2f%%.\n", index+1, hackerList[index].Login, int(hackerList[index].Value), 100*hackerList[index].Value/wholeCommentValue)
	}
	fmt.Println(content)

	fout, err := os.Create(fileUrl)
	if err != nil {
		fmt.Println(fileUrl, err)
	} else {
		fout.WriteString(string(content))
		fmt.Println("Generate ", fileUrl)
		fout.Close()
	}
}

func impactReport(issueList []github.Issue, commentList []github.Comment, fileUrl string) {
	impactMap := make(map[string]float64)
	//user who takes into consideration should at least make 'buttom_line' comments
	buttom_line := 2.0

	issueMap := make(map[string]float64)
	for index := 0; index < len(issueList); index++ {
		issue := issueList[index]
		v, ok := issueMap[issue.User.Login]
		if ok {
			v = v + 1.0
		} else {
			v = 1.0
		}
		issueMap[issue.User.Login] = v
	}
	//key: comment-login, value: counts
	commentMap := make(map[string]float64)
	for index := 0; index < len(commentList); index++ {
		comment := commentList[index]
		v, ok := commentMap[comment.User.Login]
		if ok {
			v++
		} else {
			v = 1
		}
		commentMap[comment.User.Login] = v
	}

	//key: issue-number, value: issue-login
	issueCommentMap := make(map[int]string)
	for index := 0; index < len(issueList); index++ {
		issue := issueList[index]
		//Principle 1 : drop the zero reply issue
		if issue.Comments == 0 {
			continue
		}
		//Principle 2: some people just post an issue but never replies..
		v, ok := commentMap[issue.User.Login]
		if ok == false {
			continue
		} else {
			if v < buttom_line {
				continue
			}
		}
		issueCommentMap[issueList[index].Number] = issueList[index].User.Login
		impactMap[issue.User.Login] = 0.0
	}

	for cIndex := 0; cIndex < len(commentList); cIndex++ {
		comment := commentList[cIndex]

		issueID := ospaf.GetIssueID(comment.Html_url)
		if issueID <= 0 {
			fmt.Println("Cannot get the issue ID: ", comment.Html_url)
			continue
		}
		issueLogin, ok := issueCommentMap[issueID]
		//dropped cause Principle 2
		if ok == false {
			continue
		}

		//Donnot count own impact
		if issueLogin == comment.User.Login {
			continue
		}

		_, ok = impactMap[comment.User.Login]
		if ok == false {
			continue
		}
		impactV, ok := impactMap[issueLogin]
		if ok == false {
			continue
		}
		v, _ := commentMap[comment.User.Login]
		impactMap[issueLogin] = impactV + 1.0/v
	}

	var hackerList HackerList
	wholeImpactValue := 0.0
	for key, value := range impactMap {
		var hacker Hacker
		hacker.Login = key
		hacker.Value = value
		hackerList = append(hackerList, hacker)
		wholeImpactValue += hacker.Value
	}
	sort.Sort(sort.Reverse(hackerList))

	content := fmt.Sprintf("There are %d impacts contributed by %d hackers (who has reported and commented), sort by impacts:\n", int(wholeImpactValue), len(hackerList))
	for index := 0; index < len(hackerList); index++ {
		content += fmt.Sprintf("%d. %s has %.2f impacts (%.4f per issue).\n", index+1, hackerList[index].Login, hackerList[index].Value, hackerList[index].Value/issueMap[hackerList[index].Login])
	}

	var hackerList2 HackerList
	for key, value := range impactMap {
		var hacker Hacker
		hacker.Login = key
		hacker.Value = value / issueMap[key]
		hackerList2 = append(hackerList2, hacker)
	}
	sort.Sort(sort.Reverse(hackerList2))

	content += fmt.Sprintf("\n\n-----------------------------------\nAnd sort by impacts per issue:\n")
	for index := 0; index < len(hackerList2); index++ {
		content += fmt.Sprintf("%d. %s has %.4f impacts per issue).\n", index+1, hackerList2[index].Login, hackerList2[index].Value)
	}

	fmt.Println(content)

	fout, err := os.Create(fileUrl)
	if err != nil {
		fmt.Println(fileUrl, err)
	} else {
		fout.WriteString(string(content))
		fmt.Println("Generate ", fileUrl)
		fout.Close()
	}
}

func main() {
	content, err := ospaf.ReadFile("./config.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	var config Config
	json.Unmarshal([]byte(content), &config)

	fileUrl := fmt.Sprintf("data/issue-of-%s-%s", config.Owner, config.Repo)
	content, err = ospaf.ReadFile(fileUrl)
	if err != nil {
		fmt.Println(err)
		return
	}

	var issueList []github.Issue
	json.Unmarshal([]byte(content), &issueList)

	var commentList []github.Comment
	for index := 0; index < len(issueList); index++ {
		issue := issueList[index]
		fileUrl = fmt.Sprintf("data/comment-of-issue-%s-%s-%d", config.Owner, config.Repo, issue.Number)
		content, err = ospaf.ReadFile(fileUrl)
		if err != nil {
			continue
		}

		var comments []github.Comment
		json.Unmarshal([]byte(content), &comments)
		for cIndex := 0; cIndex < len(comments); cIndex++ {
			commentList = append(commentList, comments[cIndex])
		}
	}

	ospaf.PreparePath("report", "")

	fileUrl = fmt.Sprintf("report/open-issue-%s-%s", config.Owner, config.Repo)
	openIssueReport(issueList, fileUrl)

	fileUrl = fmt.Sprintf("report/issue-comment-%s-%s", config.Owner, config.Repo)
	issueCommentReport(commentList, fileUrl)

	fileUrl = fmt.Sprintf("report/impact-comment-%s-%s", config.Owner, config.Repo)
	impactReport(issueList, commentList, fileUrl)
}
