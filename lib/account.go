package lib

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"
)

type TransactionType byte

const (
	Unknown        TransactionType = iota
	InitialBooking                 = 1
	Scooping                       = 2
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

type Account struct {
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
	Transactions  []Transaction
	Scoopings     []Transaction
}

type Transaction struct {
	Amount  int64
	Date    time.Time
	From    string
	Purpose string
	Typ     TransactionType
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

func calcGrowPer20Minutes() int64 {
	growPer20Minutes := int64(165) +
		int64(min(account.Level_1_count, 10))*int64(25) +
		int64(min(account.Level_2_count, 100))*int64(5) +
		int64(min(account.Level_3_count, 1000))*int64(1)

	return growPer20Minutes
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
