package main

import (
	ospaf "../../lib"
	"fmt"
)

func main() {
	pool, err := ospaf.InitPool()
	if err != nil {
		fmt.Println(err)
		return
	}

	level := 3
	charSet := "abcdefghijklmnopqrstuvwxyz"
	array := make([]int, level)
	for index := 0; index < level; index++ {
		array[index] = 0
	}
	array[2] = 7
	for out := false; out == false; {
		var value string
		for index := 0; index < level; index++ {
			value = fmt.Sprintf("%c", charSet[array[index]]) + value
		}
		url := fmt.Sprintf("https://api.github.com/users/%s", value)
		_, statusCode := pool.ReadURL(url, nil)
		if statusCode == -1 {
			return
		} else if statusCode != 200 {
			fmt.Println(value, statusCode)
		}

		array[0]++
		for pos := 0; pos < level; pos++ {
			if array[pos] == len(charSet) {
				array[pos] = 0
				if pos+1 < level {
					array[pos+1]++
				} else {
					out = true
					break
				}
			}
		}
	}
}
