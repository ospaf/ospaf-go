package main

import (
	ospaf "../../lib"
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func processLine(value string, pool ospaf.Pool) int {
	url := fmt.Sprintf("https://api.github.com/users/%s", value)
	_, statusCode := pool.ReadURL(url, nil)
	if statusCode == -1 {
		return -1
	} else if statusCode != 200 {
		fmt.Println(value, statusCode)
		return 1
	}
	return 0
}

func ReadLine(filePth string, pool ospaf.Pool, hookfn func(string, ospaf.Pool) int) error {
	f, err := os.Open(filePth)
	if err != nil {
		return err
	}
	defer f.Close()
	bfRd := bufio.NewReader(f)
	for {
		line, err := bfRd.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		if hookfn(strings.TrimSpace(string(line)), pool) == -1 {
			break
		}
	}
	return nil
}

func main() {
	pool, err := ospaf.InitPool()
	if err != nil {
		fmt.Println(err)
		return
	}

	ReadLine("words.txt", pool, processLine)
}
