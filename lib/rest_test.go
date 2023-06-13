package lib

import "testing"

func TestCreateAccount(t *testing.T) {
	Init("/tmp")
	res := CreateAccount("", "a4f388e4-d9a3-4b78-9f9e-3387b6ae87cc", "1cdae0a0-7896-4c9b-b103-6015a7ae1b22", "Timbuktu", "de", true)
	if res != NameMissing {
		t.Errorf("Expected to get Success but got %d", res)
	}

	res = CreateAccount("Hans", "", "1cdae0a0-7896-4c9b-b103-6015a7ae1b22", "Timbuktu", "de", true)
	if res != UuidMissing {
		t.Errorf("Expected to get Success but got %d", res)
	}

	res = CreateAccount("Hans", "a4f388e4-d9a3-4b78-9f9e-3387b6ae87cc", "", "Timbuktu", "de", true)
	if res != InviteMissing {
		t.Errorf("Expected to get Success but got %d", res)
	}

	res = CreateAccount("Hans", "a4f388e4-d9a3-4b78-9f9e-3387b6ae87cc", "1cdae0a0-7896-4c9b-b103-6015a7ae1b22", "", "de", true)
	if res != CountryMissing {
		t.Errorf("Expected to get Success but got %d", res)
	}

	res = CreateAccount("Hans", "a4f388e4-d9a3-4b78-9f9e-3387b6ae87cc", "1cdae0a0-7896-4c9b-b103-6015a7ae1b22", "Timbuktu", "de", true)
	if res != Success {
		t.Errorf("Expected to get Success but got %d", res)
	}
}
