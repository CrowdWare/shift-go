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
	acc := _account{Name: "Hans", Language: "en", Scooping: time.Date(2023, 12, 23, 9, 0, 0, 0, time.Local)}
	trans1 := _transaction{Amount: 13, Date: time.Date(2023, 3, 4, 0, 0, 0, 0, time.Local), Typ: Lmp}
	trans2 := _transaction{Amount: 15, Date: time.Date(2021, 4, 2, 0, 0, 0, 0, time.Local), Typ: Lmp}
	acc.Transactions = append(account.Transactions, trans1, trans2)
	var buffer bytes.Buffer
	acc.Name = "Bernd"
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(acc)
	if err != nil {
		t.Error(err)
	}
	acc2 := _account{}
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
	dbFile = "/tmp/shift.db"
	result := readAccount()
	expected := false
	if result != expected {
		t.Errorf("Unexpected result. Got: %v, Expected: %v", result, expected)
	}
	account = _account{}
	writeAccount()

	result = readAccount()
	expected = true
	if result != expected {
		t.Errorf("Unexpected result. Got: %v, Expected: %v", result, expected)
	}

	os.Remove("/tmp/shift.db")
}

func TestWriteReadAccount(t *testing.T) {
	dbFile = "/tmp/shift.db"
	account = _account{Name: "Hans", Language: "en", Scooping: time.Date(2023, 12, 23, 9, 0, 0, 0, time.Local)}
	trans1 := _transaction{Amount: 13, Date: time.Date(2023, 3, 4, 0, 0, 0, 0, time.Local), Typ: Lmp}
	trans2 := _transaction{Amount: 15, Date: time.Date(2021, 4, 2, 0, 0, 0, 0, time.Local), Typ: Lmp}
	account.Transactions = append(account.Transactions, trans1, trans2)
	writeAccount()

	tempAccount := _account{Name: "Hans", Language: "en", Scooping: time.Date(2023, 12, 23, 9, 0, 0, 0, time.Local)}
	tempAccount.Transactions = append(tempAccount.Transactions, trans1, trans2)

	account = _account{}
	if readAccount() != true {
		t.Error("Account could not be read")
	}
	if !reflect.DeepEqual(tempAccount, account) {
		t.Errorf("Account mismatch:\nExpected: %v\nGot: %v", tempAccount, account)
	}
	os.Remove("/tmp/shift.db")
}

func TestCheckScooping(t *testing.T) {
	dbFile = "/tmp/shift.db"
	account = _account{}
	res := checkScooping()
	if res != false {
		t.Error("Expected scooping to be off")
	}

	account.IsScooping = true
	account.Scooping = time.Now().Add(time.Hour * -20)
	res = checkScooping()
	if res == true {
		t.Error("Expected scooping to be off")
	}

	if len(account.Transactions) != 1 {
		t.Errorf("Expected transaction count to be 1 but got %d", len(account.Transactions))
	}

	account.IsScooping = true
	account.Level_1_count = 9
	account.Level_2_count = 99
	account.Level_3_count = 999
	account.Scooping = time.Now().Add(time.Hour * -20)
	res = checkScooping()
	if res == true {
		t.Error("Expected scooping to be off")
	}
	if len(account.Transactions) != 2 {
		t.Errorf("Expected transaction count to be 2 but got %d", len(account.Transactions))
	}
	balance := GetBalanceInMillis()
	if balance != 146000 {
		t.Errorf("Expected balance to be 136000 + 10000 but got %d", balance)
	}
	account.IsScooping = true
	account.Scooping = time.Now().Add(time.Hour * -19)
	res = checkScooping()
	if res == false {
		t.Error("Expected scooping to be on")
	}
	if len(account.Transactions) != 2 {
		t.Errorf("Expected transaction count to be 2 but got %d", len(account.Transactions))
	}
	os.Remove("/tmp/shift.db")
}

func TestTransactionExists(t *testing.T) {
	account = _account{}
	now := time.Now()
	addTransaction("pkey", 18, "", time.Now(), "", Lmr, "12345")
	addTransaction("", 13, "", now, "", Lmr, "12345")
	addTransaction("", 20, "", time.Now(), "", Lmr, "12345")
	res := transactionExists("pkey")
	if res == false {
		t.Error("Not implemented yet")
	}
}
