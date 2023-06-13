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
	Init("/tmp")
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
	Init("/tmp")
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

func TestAddScooping(t *testing.T) {
	Init("/tmp")

	account = _account{Level_1_count: 10, Level_2_count: 98, Level_3_count: 786}
	amount := calcGrowPer20Minutes()
	if amount != 1691 {
		t.Errorf("Expected 1691 but got %d", amount)
	}
	today := time.Now()
	for i := 0; i < 24*3; i++ {
		addScooping(calcGrowPer20Minutes(), today)
	}
	today = today.AddDate(0, 0, 1)
	addScooping(calcGrowPer20Minutes(), today)
	if len(account.Transactions) < 1 {
		t.Error("Expected 1 row count but got 0")

	} else {
		amount := account.Transactions[0].Amount
		if amount != 121 {
			t.Errorf("Expected amount to be 12 but got %d", amount)
		}
	}
	os.Remove("/tmp/shift.db")
}
