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
	From    []byte
	Purpose string
	Typ     byte
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
