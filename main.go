package main

import (
	"fmt"
	"io/ioutil"
	"regexp"

	"github.com/crowdware/shift-go/lib"
)

func readPythonConfig(filePath string) (string, string, string, error) {
	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to read file: %v", err)
	}

	re := regexp.MustCompile(`SHIFT_API_KEY\s*=\s*"(.*?)"`)
	match := re.FindStringSubmatch(string(fileContent))
	if len(match) != 2 {
		return "", "", "", fmt.Errorf("unable to find SHIFT_API_KEY in the file")
	}
	apiKey := match[1]

	re = regexp.MustCompile(`SHIFT_SECRET_KEY\s*=\s*"(.*?)"`)
	match = re.FindStringSubmatch(string(fileContent))
	if len(match) != 2 {
		return "", "", "", fmt.Errorf("unable to find SHIFT_SECRET_KEY in the file")
	}
	secretkey := match[1]

	re = regexp.MustCompile(`STORJ_ACCESS_TOKEN\s*=\s*"(.*?)"`)
	match = re.FindStringSubmatch(string(fileContent))
	if len(match) != 2 {
		return "", "", "", fmt.Errorf("unable to find STORJ_ACCESS_TOKEN in the file")
	}
	storj := match[1]

	return apiKey, secretkey, storj, nil
}

func main() {
	fmt.Println("Shift")

	configFilePath := "config.py"

	api_key, secret_key, storj_access_token, err := readPythonConfig(configFilePath)
	if err != nil {
		fmt.Println("Error reading Python config:", err)
		return
	}
	secret_key_enc := lib.Encrypt(secret_key)
	fmt.Println("const secret_key_enc = \"" + secret_key_enc + "\"")

	api_key_enc := lib.Encrypt(api_key)
	fmt.Println("const api_key_enc = \"" + api_key_enc + "\"")

	service_url_enc := lib.Encrypt("http://shift.crowdware.at:8080/")
	fmt.Println("const service_url_enc = \"" + service_url_enc + "\"")

	storjAccessTokenEnc := lib.Encrypt(storj_access_token)
	fmt.Println("const storjAccessTokenEnc = \"" + storjAccessTokenEnc + "\"")
}
