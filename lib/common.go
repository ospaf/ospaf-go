package ospafLib

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
)

func ReadFile(file_url string) (content string, err error) {
	_, err = os.Stat(file_url)
	if err != nil {
		content = fmt.Sprintf("Cannot find the file %s.", file_url)
		return content, err
	}
	file, err := os.Open(file_url)
	defer file.Close()
	if err != nil {
		content = fmt.Sprintf("Cannot open the file %s.", file_url)
		return content, err
	}
	buf := bytes.NewBufferString("")
	buf.ReadFrom(file)
	content = buf.String()

	return content, nil
}

//testStr := "https://api.github.com/repositories/36960293/issues/20/comments?page=2>; rel=\"last\", <https://api.github.com/repositories/36960293/issues/20/comments?page=1>; rel=\"first\", <https://api.github.com/repositories/36960293/issues/20/comments?page=23>; rel=\"prev\""
func GetPageMap(link string) (pageMap map[string]int) {
	pageMap = make(map[string]int)
	strSet := strings.Split(link, ",")
	for index := 0; index < len(strSet); index++ {
		re, _ := regexp.Compile("page=(\\d+)>; rel=\"(last|first|prev|next)\"")
		result := re.FindStringSubmatch(strSet[index])
		if len(result) == 3 {
			pageMap[result[2]], _ = strconv.Atoi(result[1])
		}
	}

	return pageMap
}

//testStr := "https://github.com/opencontainers/specs/pull/207#issuecomment-144513552"
func GetIssueID(link string) int {
	val := -1
	re, _ := regexp.Compile("pull/(\\d+)#issuecomment")
	result := re.FindStringSubmatch(link)
	if len(result) == 2 {
		val, _ = strconv.Atoi(result[1])
	}
	return val
}

func MD5(data string) (val string) {
	t := md5.New()
	io.WriteString(t, data)
	return fmt.Sprintf("%x", t.Sum(nil))
}

//WHen filename is null, we just want to prepare a pure directory
func PreparePath(cachename string, filename string) (realurl string) {
	var dir string
	if filename == "" {
		dir = cachename
	} else {
		realurl = path.Join(cachename, filename)
		dir = path.Dir(realurl)
	}
	p, err := os.Stat(dir)
	if err != nil {
		if !os.IsExist(err) {
			os.MkdirAll(dir, 0777)
		}
	} else {
		if p.IsDir() {
			return realurl
		} else {
			os.Remove(dir)
			os.MkdirAll(dir, 0777)
		}
	}
	return realurl
}
