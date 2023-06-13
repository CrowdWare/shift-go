package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
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
	url := decryptStringGCM(servive_url_enc) + "register"
	jsonParams := make(map[string]interface{})
	testValue := "false"
	if test {
		testValue = "true"
	}
	api_key := decryptStringGCM(api_key_enc)[:16]
	jsonParams["key"] = encryptStringGCM(api_key, true)
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
		writeAccount()
		return Success
	}
}

func SetScooping(test bool) int {
	client := http.Client{}
	url := decryptStringGCM(servive_url_enc) + "setscooping"
	jsonParams := make(map[string]interface{})
	testValue := "false"
	if test {
		testValue = "true"
	}
	api_key := decryptStringGCM(api_key_enc)[:16]
	jsonParams["key"] = encryptStringGCM(api_key, true)
	jsonParams["uuid"] = account.Uuid
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
		account.Level_1_count = int(jsonResponse["count_1"].(float64))
		account.Level_2_count = int(jsonResponse["count_2"].(float64))
		account.Level_3_count = int(jsonResponse["count_3"].(float64))
		account.Scooping = time.Now()
		writeAccount()
		return Success
	}
}

func GetMatelist(test bool) (int, []Friend) {
	emptyList := make([]Friend, 0)
	client := http.Client{}
	url := decryptStringGCM(servive_url_enc) + "matelist"
	jsonParams := make(map[string]interface{})
	testValue := "false"
	if test {
		testValue = "true"
	}
	api_key := decryptStringGCM(api_key_enc)[:16]
	jsonParams["key"] = encryptStringGCM(api_key, true)
	jsonParams["uuid"] = account.Uuid
	jsonParams["test"] = testValue

	jsonBytes, _ := json.Marshal(jsonParams)
	entity := bytes.NewBuffer(jsonBytes)

	req, _ := http.NewRequest("POST", url, entity)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", user_agent)

	resp, err := client.Do(req)
	if err != nil {
		return NetworkError, emptyList
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: %s\n", resp.Status)
		return ServiceError, emptyList
	}
	body, _ := ioutil.ReadAll(resp.Body)
	jsonResponse := make(map[string]interface{})
	json.Unmarshal(body, &jsonResponse)
	isError := jsonResponse["isError"].(bool)
	message := jsonResponse["message"].(string)
	if isError {
		log.Println("Error occured calling [register]: " + message)
		return ServiceError, emptyList
	} else {
		data := jsonResponse["data"].([]interface{})
		dataList := make([]Friend, len(data))
		for i, item := range data {
			dataObj := item.(map[string]interface{})
			name := dataObj["name"].(string)
			scooping := dataObj["scooping"].(bool)
			uuid := dataObj["uuid"].(string)
			country := dataObj["country"].(string)
			friendsCount := int(dataObj["friends_count"].(float64))
			friend := Friend{
				Name:         name,
				Scooping:     scooping,
				Uuid:         uuid,
				Country:      country,
				FriendsCount: friendsCount,
			}
			dataList[i] = friend
		}
		return Success, dataList
	}
}
