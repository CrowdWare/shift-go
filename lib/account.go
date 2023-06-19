package lib

import (
	"bytes"
	"encoding/gob"
	"log"
	"math"
	"strings"
	"time"

	"github.com/google/uuid"
)

type TransactionType byte

const (
	Unknown        TransactionType = iota
	InitialBooking                 = 1
	Scooped                        = 2
	Lmp                            = 3 // liquid micro payment
	Lmr                            = 4 // liquid micro receive
	Subtotal                       = 5
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
	Pkey    string
	Amount  int64
	Date    time.Time
	From    string
	Purpose string
	Typ     TransactionType
	Uuid    string // receiver uuid
}

func addAccount(name, _uuid, ruuid, country, language string, test bool) int {
	account = _account{
		Name:     strings.TrimSpace(name),
		Uuid:     strings.TrimSpace(_uuid),
		Ruuid:    strings.TrimSpace(ruuid),
		Country:  country,
		Language: language,
	}
	res := registerAccount(name, _uuid, ruuid, country, language, test)
	addTransaction(uuid.New().String(), initialAmount, "", time.Now(), "", InitialBooking, "")
	writeAccount()
	return res
}

func addTransaction(pkey string, amount int64, purpose string, date time.Time, from string, typ TransactionType, _uuid string) error {
	balance := GetBalanceInMillis() / 1000
	if balance+amount < 0 {
		return &BalanceError{"Amount cannot be payed out, balance to low."}
	}
	account.Transactions = append(account.Transactions, _transaction{Pkey: pkey, Amount: amount, Date: date, From: from, Purpose: purpose, Typ: typ, Uuid: _uuid})
	if len(account.Transactions) > 30 {
		// create a subtotal and delete first transaction
		account.Transactions[1].Pkey = uuid.New().String()
		account.Transactions[1].Amount += account.Transactions[0].Amount
		account.Transactions[1].Purpose = ""
		account.Transactions[1].From = ""
		account.Transactions[1].Typ = Subtotal
		account.Transactions[1].Uuid = ""
		account.Transactions = account.Transactions[1:]
	}
	writeAccount()
	return nil
}

func readAccount() bool {
	if fileExists(dbFile) {
		buffer, err := readFile(dbFile)
		if err != nil {
			if debug {
				log.Println("Error reading file " + dbFile + ", " + err.Error())
			}
			return true
		}
		decoder := gob.NewDecoder(bytes.NewReader(buffer))
		err = decoder.Decode(&account)
		if err != nil {
			if debug {
				log.Println("readAccount: " + err.Error())
			}
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
		if debug {
			log.Println(err)
		}
		return
	}
	err = writeFile(dbFile, buffer.Bytes())
	if err != nil {
		if debug {
			log.Println(err)
		}
		return
	}
}

func calcGrowPerDay() int64 {
	grow := int64(growLevel0) +
		int64(min(account.Level_1_count, 10))*growLevel1 +
		int64(min(account.Level_2_count, 100))*growLevel2 +
		int64(min(account.Level_3_count, 1000))*growLevel3
	return grow / 1000
}

func calcGrowPerDiff(duration time.Duration) int64 {
	hours := math.Min(20, duration.Hours())
	grow := float64(growLevel0)/20*hours +
		float64(min(account.Level_1_count, 10))*float64(growLevel1)/20*hours +
		float64(min(account.Level_2_count, 100))*float64(growLevel2)/20*hours +
		float64(min(account.Level_3_count, 1000))*float64(growLevel3)/20*hours
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
		addTransaction(uuid.New().String(), calcGrowPerDay(), "", time.Now(), "", Scooped, "")
		writeAccount()
	}
	return account.IsScooping
}

func calculateWorthInMillis(amount int64, transactionDate time.Time) int64 {
	currentDate := time.Now()
	daysPassed := int64(currentDate.Sub(transactionDate).Hours() / 24)
	demurrageRate := 0.27 / 100
	worth := int64(math.Pow(1-demurrageRate, float64(daysPassed)) * 1000)
	worth *= amount
	return worth
}

func transactionExists(pkey string) bool {
	for _, t := range account.Transactions {
		if t.Pkey == pkey {
			return true
		}
	}
	return false
}
