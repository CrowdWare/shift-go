package lib

import (
	"bytes"
	"encoding/gob"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestSerialize(t *testing.T) {
	acc := Account{Name: "Hans", Language: "en", Scooping: time.Date(2023, 12, 23, 9, 0, 0, 0, time.Local)}
	trans1 := Transaction{Amount: 13, Date: time.Date(2023, 3, 4, 0, 0, 0, 0, time.Local), Typ: Lmp}
	trans2 := Transaction{Amount: 15, Date: time.Date(2021, 4, 2, 0, 0, 0, 0, time.Local), Typ: Lmp}
	acc.Transactions = append(account.Transactions, trans1, trans2)
	var buffer bytes.Buffer
	acc.Name = "Bernd"
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(acc)
	if err != nil {
		t.Error(err)
	}
	acc2 := Account{}
	decoder := gob.NewDecoder(bytes.NewReader(buffer.Bytes()))
	err = decoder.Decode(&acc2)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(acc, acc2) {
		t.Errorf("Account mismatch:\nExpected: %v\nGot: %v", acc, acc2)
	}
}

func TestReadAccount(t *testing.T) {
	Init("/tmp")
	result := readAccount()
	expected := false
	if result != expected {
		t.Errorf("Unexpected result. Got: %v, Expected: %v", result, expected)
	}
	account = Account{}
	writeAccount()

	result = readAccount()
	expected = true
	if result != expected {
		t.Errorf("Unexpected result. Got: %v, Expected: %v", result, expected)
	}

	os.Remove("/tmp/shift.db")
}

func TestWriteReadAccount(t *testing.T) {
	Init("/tmp")
	account = Account{Name: "Hans", Language: "en", Scooping: time.Date(2023, 12, 23, 9, 0, 0, 0, time.Local)}
	trans1 := Transaction{Amount: 13, Date: time.Date(2023, 3, 4, 0, 0, 0, 0, time.Local), Typ: Lmp}
	trans2 := Transaction{Amount: 15, Date: time.Date(2021, 4, 2, 0, 0, 0, 0, time.Local), Typ: Lmp}
	account.Transactions = append(account.Transactions, trans1, trans2)
	writeAccount()

	tempAccount := Account{Name: "Hans", Language: "en", Scooping: time.Date(2023, 12, 23, 9, 0, 0, 0, time.Local)}
	tempAccount.Transactions = append(tempAccount.Transactions, trans1, trans2)

	account = Account{}
	if readAccount() != true {
		t.Error("Account could not be read")
	}
	if !reflect.DeepEqual(tempAccount, account) {
		t.Errorf("Account mismatch:\nExpected: %v\nGot: %v", tempAccount, account)
	}
	os.Remove("/tmp/shift.db")
}
