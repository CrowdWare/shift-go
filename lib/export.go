package lib

import (
	"math"
	"time"
)

func GetBalance() int64 {
	balance := int64(0)
	for _, t := range account.Transactions {
		balance += t.Amount
	}
	return balance
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
		writeAccount()
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
	writeAccount()
}
