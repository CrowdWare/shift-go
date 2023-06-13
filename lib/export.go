package lib

import (
	"encoding/json"
	"math"
	"time"
)

type TransactionTO struct {
	Amount  int64
	Date    int64
	From    string
	Purpose string
	Typ     int
}

func Init(filesDir string) {
	dbFile = filesDir + "/shift.db"
	account.IsScooping = false
}

func GetUuid() string {
	return account.Uuid
}

func HasJoined() bool {
	return readAccount()
}

func IsScooping() bool {
	return account.IsScooping
}

func CreateAccount(name, uuid, ruuid, country, language string) int {
	return createAccount(name, uuid, ruuid, country, language, false)
}

func GetMatelist() string {
	_, list := getMatelist(false)
	jsonData, err := json.Marshal(list)
	if err != nil {
		return ""
	}
	return string(jsonData)
}

func StartScooping() int {
	if account.IsScooping {
		return 1
	}
	res := setScooping(false)
	if res == Success {
		AddTransaction(10, "", time.Now(), "", InitialBooking)
		return Success
	}
	return 2
}

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

func GetTransactions() string {
	trans := make([]TransactionTO, 0)
	startIndex := len(account.Transactions) - 30
	if startIndex < 0 {
		startIndex = 0
	}
	for _, t := range account.Transactions[startIndex:] {
		trans = append(trans, TransactionTO{Amount: t.Amount, Purpose: t.Purpose, Date: t.Date.Unix(), From: t.From, Typ: int(t.Typ)})
	}
	jsonData, err := json.Marshal(trans)
	if err != nil {
		return ""
	}
	return string(jsonData)
}

func AddTransaction(amount int64, purpose string, date time.Time, from string, typ TransactionType) {
	account.Transactions = append(account.Transactions, _transaction{Amount: amount, Date: date, From: from, Purpose: purpose, Typ: typ})
	if len(account.Transactions) > 30 {
		// create a subtotal and delete first transaction
		account.Transactions[1].Amount += account.Transactions[0].Amount
		account.Transactions[1].Purpose = ""
		account.Transactions[1].From = ""
		account.Transactions[1].Typ = Subtotal
		account.Transactions = account.Transactions[1:]
	}
	if typ != Scooped { // don't write when Scooped
		writeAccount()
	}
}
