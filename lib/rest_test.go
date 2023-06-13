package lib

import (
	"os"
	"testing"
)

func TestCreateAccount(t *testing.T) {
	Init("/tmp")
	res := createAccount("", "a4f388e4-d9a3-4b78-9f9e-3387b6ae87cc", "1cdae0a0-7896-4c9b-b103-6015a7ae1b22", "Timbuktu", "de", true)
	if res != NameMissing {
		t.Errorf("Expected to get Success but got %d", res)
	}

	res = createAccount("Hans", "", "1cdae0a0-7896-4c9b-b103-6015a7ae1b22", "Timbuktu", "de", true)
	if res != UuidMissing {
		t.Errorf("Expected to get Success but got %d", res)
	}

	res = createAccount("Hans", "a4f388e4-d9a3-4b78-9f9e-3387b6ae87cc", "", "Timbuktu", "de", true)
	if res != InviteMissing {
		t.Errorf("Expected to get Success but got %d", res)
	}

	res = createAccount("Hans", "a4f388e4-d9a3-4b78-9f9e-3387b6ae87cc", "1cdae0a0-7896-4c9b-b103-6015a7ae1b22", "", "de", true)
	if res != CountryMissing {
		t.Errorf("Expected to get Success but got %d", res)
	}

	res = createAccount("Hans", "a4f388e4-d9a3-4b78-9f9e-3387b6ae87cc", "1cdae0a0-7896-4c9b-b103-6015a7ae1b22", "Timbuktu", "de", true)
	if res != Success {
		t.Errorf("Expected to get Success but got %d", res)
	}

	os.Remove("/tmp/shift.db")
}

func TestSetScooping(t *testing.T) {
	Init("/tmp")
	res := setScooping(true)
	if res != Success {
		t.Errorf("Expected to get Success but got %d", res)
	}
	readAccount()
	if account.Level_1_count != 9 {
		t.Errorf("Expected to get Level1 as 9 but got %d", account.Level_1_count)
	}

	if account.Level_2_count != 99 {
		t.Errorf("Expected to get Level1 as 99 but got %d", account.Level_2_count)
	}

	if account.Level_3_count != 999 {
		t.Errorf("Expected to get Level1 as 999 but got %d", account.Level_3_count)
	}
	os.Remove("/tmp/shift.db")
}

func TestGetMatelist(t *testing.T) {
	res, _ := getMatelist(true)
	if res != Success {
		t.Errorf("Expected to get Success but got %d", res)
	}
}
