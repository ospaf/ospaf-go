package main

import (
	ospaf "../"
	"fmt"
)

func main() {
	var account ospaf.Account
	account.Init("Basic", "fake001", "fake001")
	account.Load()

	fmt.Println("Account remaining: ", account.GetRemains())
	test_user := "initlove"
	url := fmt.Sprintf("https://api.github.com/users/%s", test_user)
	info, code := account.ReadURL(url, "")
	if code != 200 {
		fmt.Println(test_user, info)
	} else {
		fmt.Println(info)
	}
}
