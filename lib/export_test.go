package lib

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"log"
	"os"
	"strings"
	"testing"
	"time"
)

func TestGetTransactions(t *testing.T) {
	account = _account{}
	account.Transactions = append(account.Transactions, _transaction{Pkey: "", Amount: 34, Date: time.Date(2023, 1, 28, 4, 2, 45, 0, time.Local), Typ: InitialBooking})
	account.Transactions = append(account.Transactions, _transaction{Pkey: "", Amount: 34, Date: time.Date(2023, 1, 28, 4, 2, 45, 0, time.Local), Typ: Lmp})
	json := GetTransactions()

	if json != "[{\"Pkey\":\"\",\"Amount\":34,\"Date\":1674874965,\"From\":\"\",\"To\":\"\",\"Purpose\":\"\",\"Typ\":3},{\"Pkey\":\"\",\"Amount\":34,\"Date\":1674874965,\"From\":\"\",\"To\":\"\",\"Purpose\":\"\",\"Typ\":1}]" {
		t.Errorf("Not expected %s", json)
	}

	for i := 0; i < 30; i++ {
		account.Transactions = append(account.Transactions, _transaction{Pkey: "", Amount: 20, Date: time.Date(2023, 1, 28, 4, 2, 45, 0, time.Local), Typ: Lmp})
	}
	json = GetTransactions()
	if len(json) != 2431 {
		t.Errorf("Expected len to be 2191 but got %d", len(json))
	}
}

func TestCalculateWorth(t *testing.T) {
	amount := int64(1)
	today := time.Now()

	result := calculateWorthInMillis(amount, today)
	expectedWorth := int64(1000)
	if result != expectedWorth {
		t.Errorf("Expected worth of %d, but got %d", expectedWorth, result)
	}

	result = calculateWorthInMillis(amount, today.AddDate(0, 0, -30))
	expectedWorth = int64(922)
	if result != expectedWorth {
		t.Errorf("Expected worth of %d, but got %d", expectedWorth, result)
	}
	// nothing left after 7 years
	sevenYearsAgo := today.AddDate(0, -7*12, 0)
	expectedWorth = int64(0)
	result = calculateWorthInMillis(amount, sevenYearsAgo)
	if result != expectedWorth {
		t.Errorf("Expected worth of %d, but got %d", expectedWorth, result)
	}

	// seven years - 2 days
	sevenYearsAgo = today.AddDate(0, -7*12, 2)
	expectedWorth = int64(1)
	result = calculateWorthInMillis(amount, sevenYearsAgo)
	if result != expectedWorth {
		t.Errorf("Expected worth of %d, but got %d", expectedWorth, result)
	}

	// seven years - 2 days - double amount
	sevenYearsAgo = today.AddDate(0, -7*12, 2)
	expectedWorth = int64(2)
	result = calculateWorthInMillis(amount+1, sevenYearsAgo)
	if result != expectedWorth {
		t.Errorf("Expected worth of %d, but got %d", expectedWorth, result)
	}
}

func TestAddTransaction(t *testing.T) {
	dbFile = "/tmp/shift.db"
	account = _account{}
	addTransaction("", 10, "Purp", time.Now(), "fr", "to", InitialBooking, "")
	for i := 0; i < 31; i++ {
		addTransaction("pkey", 10, "Purp", time.Now(), "fr", "to", Scooped, "")
	}
	if len(account.Transactions) != 30 {
		t.Errorf("Expected transaction length is 30 but got %d", len(account.Transactions))
	}

	balance := GetBalanceInMillis()
	if balance != 320000 {
		t.Errorf("Expected balance is 320000 but got %d", balance)
	}

	addTransaction("", 34, "Purp", time.Now(), "ssd", "to", Lmp, "")

	balance = GetBalanceInMillis()
	if balance != 354000 {
		t.Errorf("Expected balance is 354000 but got %d", balance)
	}

	if account.Transactions[0].Typ != Subtotal {
		t.Errorf("Exception typ to be Subtotal but got %d", account.Transactions[0].Typ)
	}
	os.Remove("/tmp/shift.db")
}

func TestGetScoopedBalance(t *testing.T) {
	account.IsScooping = false
	account.Level_1_count = 9
	account.Level_2_count = 99
	account.Level_3_count = 999
	balance := GetScoopedBalance()
	if balance != 0 {
		t.Errorf("Expected balance is 0 but got %d", balance)
	}
	account.IsScooping = true
	account.Scooping = time.Now().Add(time.Hour * -3)
	balance = GetScoopedBalance()
	if balance != 20514 {
		t.Errorf("Expected balance is 20514 but got %d", balance)
	}

	account.Scooping = time.Now().Add(time.Second * -59)
	balance = GetScoopedBalance()
	if balance != 112 {
		t.Errorf("Expected balance is 112 but got %d", balance)
	}

	account.Scooping = time.Now().Add(time.Hour * -20)
	balance = GetScoopedBalance()
	if balance != 136765 {
		t.Errorf("Expected balance is 136765 but got %d", balance)
	}
}

func TestGetScoopedHours(t *testing.T) {
	account.Scooping = time.Now().Add(time.Minute * -150)
	account.IsScooping = true
	res := int64(GetScoopingHours() * 10)
	if res != 25 {
		t.Errorf("Expected 25 but got %d", res)
	}
}

func TestGetProposalQRCode(t *testing.T) {
	trans := _transaction{}
	enc := GetProposalQRCode(13, "Massage")
	plain := GetProposalFromQRCode(enc)
	err := json.Unmarshal([]byte(plain), &trans)
	if err != nil {
		log.Println(err)
	}
	if trans.Amount != -13 {
		t.Errorf("Expected amount to be -13 but got %d", trans.Amount)
	}

	if trans.From != account.Name {
		t.Errorf("Expected name to be %s but got %s", account.Name, trans.From)
	}

	if trans.Typ != Lmp {
		t.Errorf("Expected typ to be LMP but got %d", trans.Typ)
	}
}

func TestAcceptProposal(t *testing.T) {
	dbFile = "/tmp/shift.db"
	account = _account{}
	addTransaction("pkey", 10, "", time.Now(), "", "", InitialBooking, "")
	enc := GetProposalQRCode(5, "Purpose")
	plain := GetProposalFromQRCode(enc)
	if plain == "FRAUD" {
		t.Error("Expected to get a valid json code but got FRAUD")
	}
	errString := AcceptProposal()
	if errString != "ok" {
		t.Errorf("Expected to get an ok but got %s", errString)
	}
	balance := GetBalanceInMillis()
	if balance != 10000-5000 {
		t.Errorf("Expected to get a balance of 5000 but got %d", balance)
	}

	enc = GetProposalQRCode(15, "Purpose")
	plain = GetProposalFromQRCode(enc)
	if plain == "FRAUD" {
		t.Error("Expected to get a valid json code but got FRAUD")
	}
	errString = AcceptProposal()
	if errString == "ok" {
		t.Error("Expected to get an error because balance would become < 0")
	}
}

func TestGetAgreementQRCode(t *testing.T) {
	dbFile = "/tmp/shift.db"
	account = _account{}
	lastTransaction.Amount = -13
	lastTransaction.Purpose = "Purpose"
	lastTransaction.Uuid = account.Uuid
	enc := GetAgreementQRCode()
	res := GetAgreementFromQRCode(enc)
	if res != "ok" {
		t.Errorf("Expected result ok but got %s", res)
		return
	}
	if account.Transactions[0].Amount != 13 {
		t.Errorf("Expected amount to be 13 but got %d", account.Transactions[0].Amount)
	}

	if account.Transactions[0].From != account.Name {
		t.Errorf("Expected name to be %s but got %s", account.Name, account.Transactions[0].From)
	}

	if account.Transactions[0].Typ != Lmr {
		t.Errorf("Expected typ to be LMR but got %d", account.Transactions[0].Typ)
	}
	if account.Transactions[0].Purpose != "Purpose" {
		t.Errorf("Expected purpose to be Purpose but got %s", account.Transactions[0].Purpose)
	}
	os.Remove("/tmp/shift.db")
}

func TestFullTransaction(t *testing.T) {
	dbFile = "/tmp/shift.db"
	account = _account{}
	addTransaction("pkey", 20, "", time.Now(), "", "", InitialBooking, "")
	enc := GetProposalQRCode(13, "Massage")
	GetProposalFromQRCode(enc)
	AcceptProposal()
	account.Transactions[1].Pkey = "pkey2"
	enc = GetAgreementQRCode()
	GetAgreementFromQRCode(enc)
	if len(account.Transactions) != 3 {
		t.Errorf("Expected 3 transaction but got %d", len(account.Transactions))
	}
	if account.Transactions[2].Amount != 13 {
		t.Errorf("Expected amount to be 13 but got %d", account.Transactions[2].Amount)
	}
	if account.Transactions[2].Purpose != "Massage" {
		t.Errorf("Expected pursose to be Massage but got %s", account.Transactions[2].Purpose)
	}
	os.Remove("/tmp/shift.db")
}

func TestGetAgreementQRCodeForTransaction(t *testing.T) {
	now := time.Now()
	dbFile = "/tmp/shift.db"
	account = _account{}
	addTransaction("pkey", 20, "", now, "", "", InitialBooking, "")
	enc := GetProposalQRCode(13, "Massage")
	GetProposalFromQRCode(enc)
	AcceptProposal()
	enc = GetAgreementQRCode()
	res := GetAgreementFromQRCode(enc)
	enc2 := GetAgreementQRCodeForTransaction(account.Transactions[1].Pkey)
	tokens := strings.Split(enc2, "|")
	res2 := GetAgreementFromQRCode(tokens[1])

	if res != res2 {
		t.Errorf("Expected two eqal jsons but got %s and %s", res, res2)
	}
	account.Transactions[1].Typ = Scooped
	enc2 = GetAgreementQRCodeForTransaction(account.Transactions[1].Pkey)
	tokens = strings.Split(enc2, "|")
	if tokens[1] != "NOT LMP" {
		t.Errorf("Expected NOT LMP but got %s", tokens[1])
	}
	enc2 = GetAgreementQRCodeForTransaction("invalid")
	tokens = strings.Split(enc2, "|")
	if tokens[1] != "NOT FOUND" {
		t.Errorf("Expected NOT FOUND but got %s", tokens[1])
	}
	os.Remove("/tmp/shift.db")
}

func TestPeerTransfer(t *testing.T) {
	peerFile = "/tmp/peers.db"
	account = _account{}
	account.Name = "Hans"

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Error("Error generating key")
	}
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})

	peerList = []_peer{}
	peer := _peer{Name: "Hans", CryptoKey: privateKeyPEM, StorjBucket: "", StorjAccessToken: "acckey"}
	peerList = append(peerList, peer)

	code := GetPeerQRCode()
	if code != "" {
		t.Error("Could get code without bucket been set")
	}
	SetStorj("bucket", "key")
	code = GetPeerQRCode()

	res := AddPeerFromQRCode(code)
	if res != true {
		t.Error("Error getting peer from QR")
	}
	if len(peerList) != 2 {
		t.Errorf("Expected len to be 2 but got %d", len(peerList))
	}
	if peerList[1].Name != "Hans" {
		t.Errorf("Expected name to be Hans but got %s", peerList[1].Name)
	}
	publicKey := &privateKey.PublicKey
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		t.Error("Failed to encode public key:")
	}
	if !bytes.Equal(peerList[1].CryptoKey, publicKeyBytes) {
		t.Error("Keys are not equal")
	}
	os.Remove("/tmp/peers.db")
}

func TestSetStorj(t *testing.T) {
	key, err := decryptStringGCM(storj_access_token_enc, false)
	if err != nil {
		t.Error("Unable to decrypt token")
	}
	res := SetStorj("shift", key)
	if res != true {
		t.Error("Expected to get true but get false")
	}
}
