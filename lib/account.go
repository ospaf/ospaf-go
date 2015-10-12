package ospafLib

import (
	github "../github"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Account struct {
	//Basic/Oauth
	Type     string
	User     string
	Password string

	Remains int
}

//TODO: add lock to Remains

func (account *Account) Init(accountType string, accountUser string, accountPassword string) {
	account.Type = accountType
	account.User = accountUser
	account.Password = accountPassword
	account.Remains = -1
}

func (account *Account) Load() {
	url := "https://api.github.com/rate_limit"
	val, code := account.ReadURL(url, nil)
	fmt.Println("load account", val)
	if code == 200 {
		rl, ok := github.RateLimitFrom(val)
		if ok {
			account.Remains = rl.Resources.Core.Remaining
		}
	}
}

func (account *Account) GetRemains() int {
	return account.Remains
}

func (account *Account) _GetRequest(url string) (resp *http.Response, resp_body []byte, err error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	switch account.Type {
	case "Basic":
		req.SetBasicAuth(account.User, account.Password)
		break
	}
	resp, err = client.Do(req)
	if err != nil {
		return resp, resp_body, err
	}
	defer resp.Body.Close()
	resp_body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp, resp_body, err
	}

	if v, ok := resp.Header["X-Ratelimit-Remaining"]; ok {
		xRemains, err := strconv.Atoi(v[0])
		if err != nil {
			account.Remains = xRemains
			return resp, resp_body, nil
		}
	}

	if url != "https://api.github.com/rate_limit" {
		account.Remains -= 1
	}

	return resp, resp_body, nil
}

func (account *Account) ReadURL(url string, param map[string]string) (string, int) {
	var new_url string
	if param != nil {
		for key, value := range param {
			if new_url == "" {
				new_url = fmt.Sprintf("%s?%s=%s", url, key, value)
			} else {
				new_url = fmt.Sprintf("%s&%s=%s", new_url, key, value)
			}
		}
	}
	if new_url == "" {
		new_url = url
	}
	resp, resp_body, err := account._GetRequest(new_url)

	if err != nil {
		return err.Error(), -1
	}
	return string(resp_body), resp.StatusCode
}

//Return next page and the end page
func (account *Account) ReadPage(url string, param map[string]string) (body string, statusCode int, nextPage int, endPage int) {
	var new_url string
	if param != nil {
		for key, value := range param {
			if new_url == "" {
				new_url = fmt.Sprintf("%s?%s=%s", url, key, value)
			} else {
				new_url = fmt.Sprintf("%s&%s=%s", new_url, key, value)
			}
		}
	}
	if new_url == "" {
		new_url = url
	}
	fmt.Println(new_url)
	resp, resp_body, err := account._GetRequest(new_url)

	if err != nil {
		return err.Error(), -1, -1, -1
	}

	body = string(resp_body)
	statusCode = resp.StatusCode
	nextPage = -1
	endPage = -1
	if v, ok := resp.Header["Link"]; ok {
		pageMap := GetPageMap(v[0])
		if val, ok := pageMap["next"]; ok {
			nextPage = val
		}
		if val, ok := pageMap["last"]; ok {
			endPage = val
		}
	}

	return body, statusCode, nextPage, endPage
}
