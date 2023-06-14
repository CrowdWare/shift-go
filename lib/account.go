package lib

import (
	"bytes"
	"encoding/gob"
	"log"
	"math"
	"strings"
	"time"
)

type TransactionType byte

const (
	Unknown        TransactionType = iota
	InitialBooking                 = 1
	Scooped                        = 2
	Lmp                            = 3
	Subtotal                       = 4
)

type Friend struct {
	Name         string
	Scooping     bool
	Uuid         string
	Country      string
	FriendsCount int
}

type _account struct {
	Name          string
	Uuid          string
	Ruuid         string
	Language      string
	Country       string
	PrivateKey    []byte
	Scooping      time.Time
	IsScooping    bool
	Level_1_count int
	Level_2_count int
	Level_3_count int
	Transactions  []_transaction
}

type _transaction struct {
	Amount  int64
	Date    time.Time
	From    string
	Purpose string
	Typ     TransactionType
}

func addAccount(name, uuid, ruuid, country, language string, test bool) {
	account = _account{
		Name:     strings.TrimSpace(name),
		Uuid:     strings.TrimSpace(uuid),
		Ruuid:    strings.TrimSpace(ruuid),
		Country:  country,
		Language: language,
	}
	AddTransaction(10, "", time.Now(), "", InitialBooking)
	writeAccount()
	registerAccount(name, uuid, ruuid, country, language, test)
}

func readAccount() bool {
	if fileExists(dbFile) {
		buffer, err := readFile(dbFile)
		if err != nil {
			log.Fatal(err)
		}
		decoder := gob.NewDecoder(bytes.NewReader(buffer))
		err = decoder.Decode(&account)
		if err != nil {
			log.Fatal(err)
		}
		return true
	}
	return false
}

func writeAccount() {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(account)
	if err != nil {
		log.Fatal(err)
	}
	err = writeFile(dbFile, buffer.Bytes())
	if err != nil {
		log.Fatal(err)
	}
}

func calcGrowPerDay() int64 {
	grow := int64(10_000) +
		int64(min(account.Level_1_count, 10))*1800 +
		int64(min(account.Level_2_count, 100))*360 +
		int64(min(account.Level_3_count, 1000))*75
	return grow / 1000
}

func calcGrowPerDiff(duration time.Duration) int64 {
	hours := math.Min(20, duration.Hours())
	grow := float64(10_000)/20*hours +
		float64(min(account.Level_1_count, 10))*1800/20*hours +
		float64(min(account.Level_2_count, 100))*360/20*hours +
		float64(min(account.Level_3_count, 1000))*75/20*hours
	return int64(grow)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func checkScooping() bool {
	if !account.IsScooping {
		return false
	}
	diff := time.Now().Sub(account.Scooping)
	if diff.Hours() > 20 {
		account.IsScooping = false
		AddTransaction(calcGrowPerDay(), "", time.Now(), "", Scooped)
		writeAccount()
	}
	return true
}
