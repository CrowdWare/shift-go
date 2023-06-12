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
	Name         string
	Language     string
	PrivateKey   []byte
	Scooping     time.Time
	IsScooping   bool
	Transactions []Transaction
}

type Transaction struct {
	Amount  uint64
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
	account.Transactions = append(account.Transactions, Transaction{Amount: uint64(amount), Date: date, From: from, Purpose: purpose, Typ: typ})
	if len(account.Transactions) > 30 {
		// create a subtotal and delete first transaction
		account.Transactions[1].Amount += account.Transactions[0].Amount
		account.Transactions[1].Purpose = ""
		account.Transactions[1].From = ""
		account.Transactions[1].Typ = Subtotal
		account.Transactions = account.Transactions[1:]
	}
}

func GetBalance() int64 {
	balance := int64(0)
	for _, t := range account.Transactions {
		balance += int64(t.Amount)
	}
	return balance
}
