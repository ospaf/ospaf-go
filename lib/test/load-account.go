package main

import (
	ospaf "../"
	"fmt"
)

func AccountTest(account ospaf.Account) {
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
func main() {
	accounts, err := ospaf.LoadAccounts("")
	if err != nil {
		fmt.Println(err)
	} else {
		for index := 0; index < len(accounts); index++ {
			fmt.Println("\nType: ", accounts[index].Type, " Name: ", accounts[index].User)
			AccountTest(accounts[index])
		}
	}
}
