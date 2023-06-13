package lib

import (
	"os"
	"testing"
	"time"
)

func TestGetTransactions(t *testing.T) {
	account = Account{}
	account.Transactions = append(account.Transactions, Transaction{Amount: 34})
	account.Transactions = append(account.Transactions, Transaction{Amount: 34})
	account.Transactions = append(account.Transactions, Transaction{Amount: 34})
	list := GetTransactions()

	if len(list) != 3 {
		t.Errorf("Expected len of 3 but got %d", len(list))
	}
	for _, trans := range list {
		if trans.Amount != 34 {
			t.Errorf("Expected amount of 34 but got %d", trans.Amount)
		}
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
	account = Account{}
	AddTransaction(10, "Purp", time.Now(), "fr", InitialBooking)
	for i := 0; i < 31; i++ {
		AddTransaction(10, "Purp", time.Now(), "fr", Scooping)
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

func TestAddScooping(t *testing.T) {
	Init("/tmp")

	account = Account{Level_1_count: 10, Level_2_count: 98, Level_3_count: 786}
	amount := calcGrowPer20Minutes()
	if amount != 1691 {
		t.Errorf("Expected 1691 but got %d", amount)
	}
	today := time.Now()
	for i := 0; i < 24*3; i++ {
		AddScooping(calcGrowPer20Minutes(), today)
	}
	today = today.AddDate(0, 0, 1)
	AddScooping(calcGrowPer20Minutes(), today)
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
