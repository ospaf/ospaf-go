package main

import (
	github "../"
	ospaf "../../lib"
	"fmt"
	"io/ioutil"
	"net/http"
)

func gen() {
	var rl github.RateLimit
	var rlu_core github.RateLimitUnit
	var rlu_search github.RateLimitUnit
	var rlu_rate github.RateLimitUnit

	rlu_core.Limit = 10
	rlu_core.Remaining = 10

	rlu_search.Limit = 10
	rlu_search.Remaining = 10

	rlu_rate.Limit = 10
	rlu_rate.Remaining = 10

	rl.Resources.Core = rlu_core
	rl.Resources.Search = rlu_search
	rl.Rate = rlu_rate

	val := rl.Marshal()
	fmt.Println(val)
}

func readFile() {
	dataFile := "data/rate-limit.json"
	dataValue, err := ospaf.ReadFile(dataFile)
	if err == nil {
		rateLimit, valid := github.RateLimitFrom(dataValue)
		if valid {
			fmt.Println(rateLimit.Marshal())
		} else {
			fmt.Println("Cannot parse ", dataFile)
		}
	} else {
		fmt.Println("Cannot Open ", dataFile)
	}
}

func readURL() {
	url := "https://api.github.com/rate_limit"
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	resp_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	rateLimit, valid := github.RateLimitFrom(string(resp_body))
	if valid {
		fmt.Println(rateLimit.Marshal())
	} else {
		fmt.Println("Cannot parse ", resp.Body)
	}

}

func main() {
	fmt.Println("--------------gen--------------")
	gen()

	fmt.Println("\n\n------------read from data/rate-limit.json --------------")
	readFile()

	fmt.Println("\n\n------------online test --------------")
	readURL()
}
