package lib

import (
	"encoding/json"
	"log"
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
}

func GetUuid() string {
	return account.Uuid
}

func HasJoined() bool {
	return readAccount()
}

func IsScooping() bool {
	return checkScooping()
}

func CreateAccount(name, uuid, ruuid, country, language string) {
	addAccount(name, uuid, ruuid, country, language, false)
}

func GetMatelist() string {
	list := getMatelist(false)
	jsonData, err := json.Marshal(list)
	if err != nil {
		return ""
	}
	return string(jsonData)
}

func StartScooping() {
	if account.IsScooping {
		return
	}
	account.IsScooping = true
	account.Scooping = time.Now()
	account.Level_1_count = 0
	account.Level_2_count = 0
	account.Level_3_count = 0
	writeAccount()
	setScooping(false)
}

func GetBalance() int64 {
	balance := int64(0)
	for _, t := range account.Transactions {
		balance += t.Amount
	}
	return balance
}

func GetScoopedBalance() int64 {
	if account.IsScooping {
		diff := time.Now().Sub(account.Scooping)
		return calcGrowPerDiff(diff)
	}
	return 0
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
	writeAccount()
}

func GetProposalQRCode(amount int64, purpose string) string {
	trans := _transaction{Amount: amount, Purpose: purpose, Date: time.Now(), Typ: Lmp, From: account.Name}
	jsonData, err := json.Marshal(trans)
	if err != nil {
		log.Println(err)
		return ""
	}
	return encryptStringGCM(string(jsonData), false)
}

func GetTransactionFromQRCode(enc string) string {
	json := decryptStringGCM(enc)
	return string(json)
}
