package ospafLib

import (
	"encoding/json"
)

//TODO: use /etc or /home?
func LoadAccounts(confIn string) ([]Account, error) {
	var accounts []Account
	var conf string
	if confIn == "" {
		conf = "accounts.json"
	} else {
		conf = confIn
	}
	content, err := ReadFile(conf)
	if err == nil {
		err = json.Unmarshal([]byte(content), &accounts)
	}

	return accounts, err
}
