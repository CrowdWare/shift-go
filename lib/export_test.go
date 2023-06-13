package lib

import (
	"testing"
	"time"
)

func TestGetTransactions(t *testing.T) {
	account = _account{}
	account.Transactions = append(account.Transactions, _transaction{Amount: 34, Date: time.Date(2023, 1, 28, 4, 2, 45, 0, time.Local), Typ: Lmp})
	account.Transactions = append(account.Transactions, _transaction{Amount: 34, Date: time.Date(2023, 1, 28, 4, 2, 45, 0, time.Local), Typ: InitialBooking})
	json := GetTransactions()

	if json != "[{\"Amount\":34,\"Date\":1674874965,\"From\":\"\",\"Purpose\":\"\",\"Typ\":3},{\"Amount\":34,\"Date\":1674874965,\"From\":\"\",\"Purpose\":\"\",\"Typ\":1}]" {
		t.Errorf("Not expected %s", json)
	}
}

func TestCalculateWorth(t *testing.T) {
	amount := int64(1)
	today := time.Now()

	result := CalculateWorthInMillis(amount, today)
	expectedWorth := int64(1000)
	if result != expectedWorth {
		t.Errorf("Expected worth of %d, but got %d", expectedWorth, result)
	}

	result = CalculateWorthInMillis(amount, today.AddDate(0, 0, -30))
	expectedWorth = int64(922)
	if result != expectedWorth {
		t.Errorf("Expected worth of %d, but got %d", expectedWorth, result)
	}
	// nothing left after 7 years
	sevenYearsAgo := today.AddDate(0, -7*12, 0)
	expectedWorth = int64(0)
	result = CalculateWorthInMillis(amount, sevenYearsAgo)
	if result != expectedWorth {
		t.Errorf("Expected worth of %d, but got %d", expectedWorth, result)
	}

	// seven years - 2 days
	sevenYearsAgo = today.AddDate(0, -7*12, 2)
	expectedWorth = int64(1)
	result = CalculateWorthInMillis(amount, sevenYearsAgo)
	if result != expectedWorth {
		t.Errorf("Expected worth of %d, but got %d", expectedWorth, result)
	}

	// seven years - 2 days - double amount
	sevenYearsAgo = today.AddDate(0, -7*12, 2)
	expectedWorth = int64(2)
	result = CalculateWorthInMillis(amount+1, sevenYearsAgo)
	if result != expectedWorth {
		t.Errorf("Expected worth of %d, but got %d", expectedWorth, result)
	}
}

func TestAddTransaction(t *testing.T) {
	Init("/tmp")
	account = _account{}
	AddTransaction(10, "Purp", time.Now(), "fr", InitialBooking)
	for i := 0; i < 31; i++ {
		AddTransaction(10, "Purp", time.Now(), "fr", Scooped)
	}
	if len(account.Transactions) != 30 {
		t.Errorf("Expected transaction length is 30 but got %d", len(account.Transactions))
	}

	balance := GetBalance()
	if balance != 320 {
		t.Errorf("Expected balance is 320 but got %d", balance)
	}

	AddTransaction(34, "Purp", time.Now(), "ssd", Lmp)

	balance = GetBalance()
	if balance != 354 {
		t.Errorf("Expected balance is 354 but got %d", balance)
	}

	if account.Transactions[0].Typ != Subtotal {
		t.Errorf("Exception typ to be Subtotal but got %d", account.Transactions[0].Typ)
	}
}
