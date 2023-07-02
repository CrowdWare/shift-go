/****************************************************************************
 * Copyright (C) 2023 CrowdWare
 *
 * This file is part of SHIFT.
 *
 *  SHIFT is free software: you can redistribute it and/or modify
 *  it under the terms of the GNU General Public License as published by
 *  the Free Software Foundation, either version 3 of the License, or
 *  (at your option) any later version.
 *
 *  SHIFT is distributed in the hope that it will be useful,
 *  but WITHOUT ANY WARRANTY without even the implied warranty of
 *  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *  GNU General Public License for more details.
 *
 *  You should have received a copy of the GNU General Public License
 *  along with SHIFT.  If not, see <http://www.gnu.org/licenses/>.
 *
 ****************************************************************************/
package lib

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
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
	if !useWebService {
		return 0
	}

	// makes no sense to register an account without invite code
	if ruuid == "" {
		return 1
	}

	client := http.Client{}

	url, err := decryptStringGCM(service_url_enc, false)
	if err != nil {
		if debug {
			log.Println("error decrypting service url: " + err.Error())
		}
		return 6
	}

	jsonParams := make(map[string]interface{})
	testValue := "false"
	if test {
		testValue = "true"
	}
	api_key, err := decryptStringGCM(api_key_enc, false)
	if err != nil {
		if debug {
			log.Println("error decrypting api key: " + err.Error())
		}
		return 5
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
		if debug {
			log.Println("Error occured calling [register]: " + err.Error())
		}
		return 2
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		if debug {
			log.Printf("Error: %s\n", resp.Status)
		}
		return 3
	}
	body, _ := ioutil.ReadAll(resp.Body)
	jsonResponse := make(map[string]interface{})
	json.Unmarshal(body, &jsonResponse)
	isError := jsonResponse["isError"].(bool)
	message := jsonResponse["message"].(string)

	if isError {
		if debug {
			log.Println("Error occured calling register: " + message)
		}
		return 4
	}

	return 0
}

func setScooping(test bool) int {
	if !useWebService {
		account.IsScooping = true
		account.Scooping = time.Now()
		account.Level_1_count = 0
		account.Level_2_count = 0
		account.Level_3_count = 0
		writeAccount()
		return 0
	}

	client := http.Client{}

	url, err := decryptStringGCM(service_url_enc, false)
	if err != nil {
		if debug {
			log.Println("Error decrypting service url:" + err.Error())
		}
		return 1
	}

	jsonParams := make(map[string]interface{})
	testValue := "false"
	if test {
		testValue = "true"
	}

	api_key, err := decryptStringGCM(api_key_enc, false)
	if err != nil {
		if debug {
			log.Println("Error decrypting api key: " + err.Error())
		}
		return 2
	}
	timeString := time.Now().Format("2006-01-02 15:04:05")
	jsonParams["key"] = encryptStringGCM(api_key[:16], true)
	jsonParams["uuid"] = account.Uuid
	jsonParams["test"] = testValue
	jsonParams["time"] = timeString

	jsonBytes, _ := json.Marshal(jsonParams)
	entity := bytes.NewBuffer(jsonBytes)

	req, _ := http.NewRequest("POST", url+"setscooping", entity)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", user_agent)

	resp, err := client.Do(req)
	if err != nil {
		return 3
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		if debug {
			log.Printf("Error: %s\n", resp.Status)
		}
		return 4
	}
	body, _ := ioutil.ReadAll(resp.Body)
	jsonResponse := make(map[string]interface{})
	json.Unmarshal(body, &jsonResponse)
	isError := jsonResponse["isError"].(bool)
	message := jsonResponse["message"].(string)
	if isError {
		if debug {
			log.Println("Error occured calling setscooping: " + message)
		}
		account.IsScooping = false
		writeAccount()
		return 5
	} else {
		client_key, err := decryptStringGCM(client_key_enc, false)
		if err != nil {
			if debug {
				log.Println("Error decrypting client key: " + err.Error())
			}
			return 6
		}
		key_enc := jsonResponse["key"].(string)
		key, err := decryptStringGCM(key_enc, true)
		if err != nil {
			if debug {
				log.Println("Error decrypting response: " + err.Error())
			}
			return 7
		}
		if key != client_key+timeString {
			if debug {
				log.Println("Error response not valid: " + key)
			}
			return 8
		}
		account.IsScooping = true
		account.Scooping = time.Now()
		account.Level_1_count = int(jsonResponse["count_1"].(float64))
		account.Level_2_count = int(jsonResponse["count_2"].(float64))
		account.Level_3_count = int(jsonResponse["count_3"].(float64))
		writeAccount()
		return 0
	}
}

func getMatelist(test bool) []Friend {
	emptyList := make([]Friend, 0)
	if !useWebService {
		return emptyList
	}

	client := http.Client{}
	url, err := decryptStringGCM(service_url_enc, false)
	if err != nil {
		if debug {
			log.Println("Error decrypting service url: " + err.Error())
		}
		return emptyList
	}

	jsonParams := make(map[string]interface{})
	testValue := "false"
	if test {
		testValue = "true"
	}
	api_key, err := decryptStringGCM(api_key_enc, false)
	if err != nil {
		if debug {
			log.Println("Error decrypting api key: " + err.Error())
		}
		return emptyList
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
		if debug {
			log.Println("Error posting matelist: " + err.Error())
		}
		return emptyList
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		if debug {
			log.Printf("Error: %s\n", resp.Status)
		}
		return emptyList
	}
	body, _ := ioutil.ReadAll(resp.Body)
	jsonResponse := make(map[string]interface{})
	json.Unmarshal(body, &jsonResponse)
	isError := jsonResponse["isError"].(bool)
	message := jsonResponse["message"].(string)
	if isError {
		if debug {
			log.Println("Error occured calling register: " + message)
		}
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
