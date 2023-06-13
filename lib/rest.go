package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

/*	Make the life of the hackers a bit harder ;-)
	Encrypt Keys before deploy as follows
	api_key_enc := lib.EncryptStringGCM(api_key_16_bytes_same_as_on_webservice + "9878623447645687")
	servive_url_enc := lib.EncryptStringGCM("http://shift.crowdware.at:8080/", false)
*/

const servive_url_enc = "2264fb60799a2cb14026bdf896aa0091da58d65d76a8694de32edd1fffb9ae9566d39cd48d8dcfe5397ebcae2b183e2b8c94ee49b5fd185ff1318c"
const api_key_enc = "1dd85261864261b7182f43d6e7a65691d20ae5382941b60d0b8c6bcbc7d5345e473859d17611c48a7923dba552d5032a46997634b025341cd6c0eeff"
const user_agent = "Shift 1.0"

type RestResult byte

const (
	Success        = iota
	NameMissing    = 1
	InviteMissing  = 2
	CountryMissing = 3
	UuidMissing    = 4
	NetworkError   = 5
	ServiceError   = 6
)

func CreateAccount(
	name string,
	uuid string,
	ruuid string,
	country string,
	language string,
	test bool,
) int {
	if name == "" {
		return NameMissing
	} else if uuid == "" {
		return UuidMissing
	} else if ruuid == "" {
		return InviteMissing
	} else if country == "" {
		return CountryMissing
	}

	account = Account{
		Name:     strings.TrimSpace(name),
		Uuid:     strings.TrimSpace(uuid),
		Ruuid:    strings.TrimSpace(ruuid),
		Country:  country,
		Language: language,
	}

	client := http.Client{}
	url := DecryptStringGCM(servive_url_enc) + "register"
	jsonParams := make(map[string]interface{})
	testValue := "false"
	if test {
		testValue = "true"
	}
	api_key := DecryptStringGCM(api_key_enc)[:16]
	jsonParams["key"] = EncryptStringGCM(api_key, true)
	jsonParams["name"] = account.Name
	jsonParams["uuid"] = account.Uuid
	jsonParams["ruuid"] = account.Ruuid
	jsonParams["country"] = account.Country
	jsonParams["language"] = account.Language
	jsonParams["test"] = testValue

	jsonBytes, _ := json.Marshal(jsonParams)
	entity := bytes.NewBuffer(jsonBytes)

	req, _ := http.NewRequest("POST", url, entity)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", user_agent)

	resp, err := client.Do(req)
	if err != nil {
		return NetworkError
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: %s\n", resp.Status)
		return ServiceError
	}
	body, _ := ioutil.ReadAll(resp.Body)
	jsonResponse := make(map[string]interface{})
	json.Unmarshal(body, &jsonResponse)
	isError := jsonResponse["isError"].(bool)
	message := jsonResponse["message"].(string)

	if isError {
		log.Println("Error occured calling [register]: " + message)
		return ServiceError
	} else {
		WriteAccount()
		return Success
	}
}
