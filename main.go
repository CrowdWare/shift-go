package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"time"

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

	re = regexp.MustCompile(`SHIFT_CLIENT_KEY\s*=\s*"(.*?)"`)
	match = re.FindStringSubmatch(string(fileContent))
	if len(match) != 2 {
		return "", "", "", fmt.Errorf("unable to find SHIFT_CLIENT_KEY in the file")
	}
	clientKey := match[1]

	re = regexp.MustCompile(`SHIFT_SECRET_KEY\s*=\s*"(.*?)"`)
	match = re.FindStringSubmatch(string(fileContent))
	if len(match) != 2 {
		return "", "", "", fmt.Errorf("unable to find SHIFT_SECRET_KEY in the file")
	}
	secretkey := match[1]

	return apiKey, secretkey, clientKey, nil
}

func generateRandomNumber(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min+1) + min
}

func writeVars() {
	file, err := os.Create("./lib/crypto_vars.go.temp")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	text := "package lib\n\n"
	text += "const debug = false\n"
	text += "const initial_amount = 1\n\n"
	text += "const var1 = " + strconv.Itoa(generateRandomNumber(1000, 9999)) + "\n"
	text += "const var2 = " + strconv.Itoa(generateRandomNumber(1000, 9999)) + "\n"
	text += "const var3 = " + strconv.Itoa(generateRandomNumber(10, 99)) + "\n"
	text += "const var4 = " + strconv.Itoa(generateRandomNumber(100000, 999999)) + "\n"
	text += "const var5 = " + strconv.Itoa(generateRandomNumber(100000, 999999)) + "\n"
	text += "\n"
	text += "const secret_key_enc = \"\"\n"
	text += "const api_key_enc = \"\"\n"
	text += "const client_key_enc = \"\"\n"
	text += "const service_url_enc = \"\"\n"

	_, err = writer.WriteString(text)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
	err = writer.Flush()
	if err != nil {
		fmt.Println("Error flushing writer:", err)
		return
	}
}

func secret() {
	infile, err := os.Open("./lib/crypto_vars.go")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer infile.Close()

	scanner := bufio.NewScanner(infile)
	text := ""
	numLines := 10
	lineCount := 0
	for scanner.Scan() {
		line := scanner.Text()
		text += line + "\n"
		lineCount++
		if lineCount == numLines {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}

	file, err := os.Create("./lib/crypto_vars.go.temp")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	text += "\n"

	api_key, secret_key, client_key, err := readPythonConfig("shift_keys.py")
	if err != nil {
		fmt.Println("Error reading Python config:", err)
		return
	}
	secret_key_enc := lib.Encrypt(secret_key)
	text += "const secret_key_enc = \"" + secret_key_enc + "\"\n"

	// there was a reason that I used this extra bytes
	api_key_enc := lib.Encrypt(api_key + "8764398347362489")
	text += "const api_key_enc = \"" + api_key_enc + "\"\n"

	client_key_enc := lib.Encrypt(client_key)
	text += "const client_key_enc = \"" + client_key_enc + "\"\n"

	service_url_enc := lib.Encrypt("http://shift.crowdware.at:8080/")
	text += "const service_url_enc = \"" + service_url_enc + "\"\n"

	_, err = writer.WriteString(text)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
	err = writer.Flush()
	if err != nil {
		fmt.Println("Error flushing writer:", err)
		return
	}
}

/*
**	Initialize secret keys.
**	First exchange variable1...variable5 in crypto
 */
func main() {
	if len(os.Args) > 1 {
		if os.Args[1] == "vars" {
			writeVars()
		} else if os.Args[1] == "secret" {
			secret()
		} else if os.Args[1] == "test" {
			test()
		} else {
			fmt.Println("Unknown argument:" + os.Args[1])
		}
	} else {
		fmt.Println("Shift")
		fmt.Println("Usage: go run . <arg> | where args are: vars | secret")
	}
}

func test() {

}
