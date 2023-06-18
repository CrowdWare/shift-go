package lib

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

const user_agent = "Shift 1.0"

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
	url, err := decryptStringGCM(service_url_enc)
	if err != nil {
		log.Println("error with decryptStringGCM")
		log.Fatal(err)
	}
	log.Println("url: " + url)

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
		log.Printf("Error: %s\n", resp.Status)
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
	url, err := decryptStringGCM(service_url_enc)
	if err != nil {
		log.Fatal(err)
	}
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
		log.Printf("Error: %s\n", resp.Status)
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
	url, err := decryptStringGCM(service_url_enc)
	if err != nil {
		log.Fatal(err)
	}

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
		log.Printf("Error: %s\n", resp.Status)
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
