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
	Scoopings     []_transaction
}

type _transaction struct {
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

func addScooping(amount int64, date time.Time) {
	// when the last scooping has been added yesterday, then sum up the scooping and create a new transaction
	len := len(account.Scoopings)
	if len > 0 && account.Scoopings[len-1].Date.Day() != date.Day() {
		milliLiter := int64(0)
		for _, t := range account.Scoopings {
			milliLiter += t.Amount
		}
		AddTransaction(milliLiter/1000, "", date, "", Scooped)
		account.Scoopings = make([]_transaction, 0)
	}
	account.Scoopings = append(account.Scoopings, _transaction{Amount: amount, Date: date})
	writeAccount()
}
