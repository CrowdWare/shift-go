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

func readPythonConfig(filePath string) (string, string, string, string, error) {
	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", "", "", "", fmt.Errorf("failed to read file: %v", err)
	}

	re := regexp.MustCompile(`SHIFT_API_KEY\s*=\s*"(.*?)"`)
	match := re.FindStringSubmatch(string(fileContent))
	if len(match) != 2 {
		return "", "", "", "", fmt.Errorf("unable to find SHIFT_API_KEY in the file")
	}
	apiKey := match[1]

	re = regexp.MustCompile(`SHIFT_CLIENT_KEY\s*=\s*"(.*?)"`)
	match = re.FindStringSubmatch(string(fileContent))
	if len(match) != 2 {
		return "", "", "", "", fmt.Errorf("unable to find SHIFT_CLIENT_KEY in the file")
	}
	clientKey := match[1]

	re = regexp.MustCompile(`SHIFT_SECRET_KEY\s*=\s*"(.*?)"`)
	match = re.FindStringSubmatch(string(fileContent))
	if len(match) != 2 {
		return "", "", "", "", fmt.Errorf("unable to find SHIFT_SECRET_KEY in the file")
	}
	secretkey := match[1]

	re = regexp.MustCompile(`STORJ_ACCESS_TOKEN\s*=\s*"(.*?)"`)
	match = re.FindStringSubmatch(string(fileContent))
	if len(match) != 2 {
		return "", "", "", "", fmt.Errorf("unable to find STORJ_ACCESS_TOKEN in the file")
	}
	storj := match[1]

	return apiKey, secretkey, storj, clientKey, nil
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
	text += "const var3 = " + strconv.Itoa(generateRandomNumber(1000, 9999)) + "\n"
	text += "const var4 = " + strconv.Itoa(generateRandomNumber(100000, 999999)) + "\n"
	text += "const var5 = " + strconv.Itoa(generateRandomNumber(100000, 999999)) + "\n"
	text += "\n"
	text += "const secret_key_enc = \"\"\n"
	text += "const api_key_enc = \"\"\n"
	text += "const client_key_enc = \"\"\n"
	text += "const service_url_enc = \"\"\n"
	text += "const storj_access_token_enc = \"\"\n"

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

	api_key, secret_key, storj_access_token, client_key, err := readPythonConfig("shift_keys.py")
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

	storj_access_token_enc := lib.Encrypt(storj_access_token)
	text += "const storj_access_token_enc = \"" + storj_access_token_enc + "\"\n"

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
	plain, err := lib.Decrypt("83198d289420abbff5bbdbba84353bcea240ff41aaca651ef1bb2f0e2cbac211872a32b3e45899e46679a9c655d1f1114573839dcc3e25ae1e1fbdc98598a752a7eec48eaabc1179dfe6cfbe7899b58824a1f0399769a5298d255363d72fff39e0cb7c8a7d7e9dbf3470a95b76e71ba05feaa554df6cefaa1aa6dcb401586f3f0557e54171de06f6324e5d058412426adc3afbcdc2d0504655006901ab67bd3647b7e5e0e66f1c32962ff9b0605db4d7a6e6329c9680f0e11a97c5be4a660761033b4d0104c767e6ba265dbe8b1a24a9c8b67c6d163739ffa6579b9bc561e8ed03f6c19439ce4771b21158c76bdc7acd070eb6ae708e0ad0988f4cb3424741ac0f4222630d35c1bcd9641dd2afd1645c56f100727d5c1f6f006617f5d4a2bd5701a8f20cba3c6ad98cdd57ad0f71224c18ce2c6003ca15d001bcd6be07e6ef498abc0874fd1a8b7836097dd077ce4e6008f5593127db822f")
	if err != nil {
		fmt.Println("Decrpyt error: " + err.Error())
		return
	}
	fmt.Println(plain == "1GW7L5Hab3vR4twJARK4mMuatA2D319NyYboQXnRQU9JcLDj2BEwwtiZ5whRtwDV4KRPvsfV4HcSjq9DutvF2NLr6yMgij6N6debnCzeLEfPZJds2uLtj4PcQHPXUyzqStdxwTAZrMDJX4RQcvdpqAtbRUVxtbrkg7hRCrjgwTFNCAoATvfeeoXacMkUBMSxpNXLfp3NYWk9KjGgbRC9SkFHDurkrHg8aVs1mMs2vRqW2Y1mcHbpzYthWJxfJB1sQP1shfRyCUZxTY4okb5gnZH3tSSyCPSsSkbLh6KSYnVrb2bqRAr1AgvfQVaB")

	plain, err = lib.Decrypt("1fb5ffeb0a0dc598f8383d8b2691b576dcaaa253518dfb0a6b47465c311d6eceb5581d5a6440935b7667837f3ea55e3b275fecb43459d5438fa208")
	if err != nil {
		fmt.Println("Decrpyt error: " + err.Error())
		return
	}
	fmt.Println(plain == "http://shift.crowdware.at:8080/")
}
