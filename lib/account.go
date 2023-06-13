package lib

import (
	"bytes"
	"encoding/gob"
	"log"
	"math"
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

func ReadAccount() bool {
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

func WriteAccount() {
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

func CalculateWorthInMillis(amount int64, transactionDate time.Time) int64 {
	currentDate := time.Now()
	daysPassed := int64(currentDate.Sub(transactionDate).Hours() / 24)
	demurrageRate := 0.27 / 100
	worth := int64(math.Pow(1-demurrageRate, float64(daysPassed)) * 1000)
	worth *= amount
	return worth
}

func GetTransactions() []Transaction {
	trans := make([]Transaction, 0)
	startIndex := len(account.Transactions) - 30
	if startIndex < 0 {
		startIndex = 0
	}
	for _, t := range account.Transactions[startIndex:] {
		trans = append(trans, t)
	}
	return trans
}

func AddTransaction(amount int64, purpose string, date time.Time, from string, typ TransactionType) {
	account.Transactions = append(account.Transactions, Transaction{Amount: amount, Date: date, From: from, Purpose: purpose, Typ: typ})
	if len(account.Transactions) > 30 {
		// create a subtotal and delete first transaction
		account.Transactions[1].Amount += account.Transactions[0].Amount
		account.Transactions[1].Purpose = ""
		account.Transactions[1].From = ""
		account.Transactions[1].Typ = Subtotal
		account.Transactions = account.Transactions[1:]
	}
	if typ != Scooping { // don't write when Scooping
		WriteAccount()
	}
}

func AddScooping(amount int64, date time.Time) {
	// when the last scooping has been added yesterday, then sum up the scooping and create a new transaction
	len := len(account.Scoopings)
	if len > 0 && account.Scoopings[len-1].Date.Day() != date.Day() {
		milliLiter := int64(0)
		for _, t := range account.Scoopings {
			milliLiter += t.Amount
		}
		AddTransaction(milliLiter/1000, "", date, "", Scooping)
		account.Scoopings = make([]Transaction, 0)
	}
	account.Scoopings = append(account.Scoopings, Transaction{Amount: amount, Date: date})
	WriteAccount()
}

func GetBalance() int64 {
	balance := int64(0)
	for _, t := range account.Transactions {
		balance += t.Amount
	}
	return balance
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
