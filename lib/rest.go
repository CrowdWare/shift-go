package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

/*	Make the life of the hackers a bit harder ;-)
	Encrypt Keys before deploy as follows
	api_key_enc := lib.EncryptStringGCM(api_key_16_bytes_same_as_on_webservice + "9878623447645687")
	servive_url_enc := lib.EncryptStringGCM("http://shift.crowdware.at:8080/", false)
*/

const servive_url_enc = "2264fb60799a2cb14026bdf896aa0091da58d65d76a8694de32edd1fffb9ae9566d39cd48d8dcfe5397ebcae2b183e2b8c94ee49b5fd185ff1318c"
const api_key_enc = "1dd85261864261b7182f43d6e7a65691d20ae5382941b60d0b8c6bcbc7d5345e473859d17611c48a7923dba552d5032a46997634b025341cd6c0eeff"
const user_agent = "Shift 1.0"
const url = "http://128.140.48.116:8080/"

func registerAccount(
	name string,
	uuid string,
	ruuid string,
	country string,
	language string,
	test bool,
) int {
	// makes no sense to register an account without invite code
	if ruuid == "" {
		return 1
	}
	client := http.Client{}
	//url, err := decryptStringGCM(servive_url_enc)
	//if err != nil {
	//	log.Fatal(err)
	//}

	jsonParams := make(map[string]interface{})
	testValue := "false"
	if test {
		testValue = "true"
	}
	api_key, err := decryptStringGCM(api_key_enc)
	if err != nil {
		log.Fatal(err)
	}

	jsonParams["key"] = encryptStringGCM(api_key[:16], true)
	jsonParams["name"] = account.Name
	jsonParams["uuid"] = account.Uuid
	jsonParams["ruuid"] = account.Ruuid
	jsonParams["country"] = account.Country
	jsonParams["language"] = account.Language
	jsonParams["test"] = testValue

	jsonBytes, _ := json.Marshal(jsonParams)
	entity := bytes.NewBuffer(jsonBytes)

	req, _ := http.NewRequest("POST", url+"register", entity)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", user_agent)

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error occured calling [register]: " + err.Error())
		return 2
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: %s\n", resp.Status)
		return 3
	}
	body, _ := ioutil.ReadAll(resp.Body)
	jsonResponse := make(map[string]interface{})
	json.Unmarshal(body, &jsonResponse)
	isError := jsonResponse["isError"].(bool)
	message := jsonResponse["message"].(string)

	if isError {
		log.Println("Error occured calling [register]: " + message)
	}
	return 0
}

func setScooping(test bool) int {
	client := http.Client{}
	//url, err := decryptStringGCM(servive_url_enc)
	//if err != nil {
	//	log.Fatal(err)
	//	}
	jsonParams := make(map[string]interface{})
	testValue := "false"
	if test {
		testValue = "true"
	}
	api_key, err := decryptStringGCM(api_key_enc)
	if err != nil {
		log.Fatal(err)
	}

	jsonParams["key"] = encryptStringGCM(api_key[:16], true)
	jsonParams["uuid"] = account.Uuid
	jsonParams["test"] = testValue

	jsonBytes, _ := json.Marshal(jsonParams)
	entity := bytes.NewBuffer(jsonBytes)

	req, _ := http.NewRequest("POST", url+"setscooping", entity)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", user_agent)

	resp, err := client.Do(req)
	if err != nil {
		return 1
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: %s\n", resp.Status)
		return 2
	}
	body, _ := ioutil.ReadAll(resp.Body)
	jsonResponse := make(map[string]interface{})
	json.Unmarshal(body, &jsonResponse)
	isError := jsonResponse["isError"].(bool)
	message := jsonResponse["message"].(string)
	if isError {
		log.Println("Error occured calling [register]: " + message)
		return 3
	} else {
		account.Level_1_count = int(jsonResponse["count_1"].(float64))
		account.Level_2_count = int(jsonResponse["count_2"].(float64))
		account.Level_3_count = int(jsonResponse["count_3"].(float64))
		writeAccount()
		return 0
	}
}

func getMatelist(test bool) []Friend {
	emptyList := make([]Friend, 0)
	client := http.Client{}
	//url, err := decryptStringGCM(servive_url_enc)
	//if err != nil {
	//	log.Fatal(err)
	//}

	jsonParams := make(map[string]interface{})
	testValue := "false"
	if test {
		testValue = "true"
	}
	api_key, err := decryptStringGCM(api_key_enc)
	if err != nil {
		log.Fatal(err)
	}

	jsonParams["key"] = encryptStringGCM(api_key[:16], true)
	jsonParams["uuid"] = account.Uuid
	jsonParams["test"] = testValue

	jsonBytes, _ := json.Marshal(jsonParams)
	entity := bytes.NewBuffer(jsonBytes)

	req, _ := http.NewRequest("POST", url+"matelist", entity)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", user_agent)

	resp, err := client.Do(req)
	if err != nil {
		return emptyList
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: %s\n", resp.Status)
		return emptyList
	}
	body, _ := ioutil.ReadAll(resp.Body)
	jsonResponse := make(map[string]interface{})
	json.Unmarshal(body, &jsonResponse)
	isError := jsonResponse["isError"].(bool)
	message := jsonResponse["message"].(string)
	if isError {
		log.Println("Error occured calling [register]: " + message)
		return emptyList
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
		return dataList
	}
}
