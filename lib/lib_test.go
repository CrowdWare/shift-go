package lib

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"errors"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestReadAccount(t *testing.T) {
	Init("/tmp")
	result := ReadAccount()
	expected := false
	if result != expected {
		t.Errorf("Unexpected result. Got: %v, Expected: %v", result, expected)
	}
	account = Account{}
	WriteAccount()

	result = ReadAccount()
	expected = true
	if result != expected {
		t.Errorf("Unexpected result. Got: %v, Expected: %v", result, expected)
	}

	os.Remove("/tmp/shift.db")
}

func TestEncryptDecrypt(t *testing.T) {
	Init("")
	teststring := "The quick brown fox"
	enc := EncryptStringGCM(teststring)
	result := DecryptStringGCM(enc)
	expected := teststring
	if result != expected || enc == teststring {
		t.Errorf("Unexpected result. Got: %s, Expected: %s", result, expected)
	}
}

func TestEncryptAndDecryptBytesGCM(t *testing.T) {
	plaintext := []byte("Hello, World!")
	Init("")
	// Encrypt the plaintext
	ciphertext, nonce, err := EncryptBytesGCM(plaintext)
	if err != nil {
		t.Errorf("Encryption error: %v", err)
	}

	// Decrypt the ciphertext
	decrypted, err := DecryptBytesGCM(ciphertext, nonce)
	if err != nil {
		t.Errorf("Decryption error: %v", err)
	}

	// Check if the decrypted plaintext matches the original plaintext
	if !bytes.Equal(decrypted, plaintext) {
		t.Errorf("Decrypted plaintext does not match the original plaintext")
	}
}

func TestEncryptAndDecryptBytesGCMWithBinaryData(t *testing.T) {
	plaintext := []byte{0x01, 0x02, 0x03, 0x04, 0x05}
	Init("")
	// Encrypt the plaintext
	ciphertext, nonce, err := EncryptBytesGCM(plaintext)
	if err != nil {
		t.Errorf("Encryption error: %v", err)
	}

	// Decrypt the ciphertext
	decrypted, err := DecryptBytesGCM(ciphertext, nonce)
	if err != nil {
		t.Errorf("Decryption error: %v", err)
	}

	// Check if the decrypted plaintext matches the original plaintext
	if !bytes.Equal(decrypted, plaintext) {
		t.Errorf("Decrypted plaintext does not match the original plaintext")
	}
}

func TestEncryptAndDecryptBytesGCMWithHexInput(t *testing.T) {
	plaintext := []byte("Hello, World!")
	Init("")
	// Encrypt the plaintext
	ciphertext, nonce, err := EncryptBytesGCM(plaintext)
	if err != nil {
		t.Errorf("Encryption error: %v", err)
	}

	// Encode ciphertext and nonce to hexadecimal strings
	ciphertextHex := hex.EncodeToString(ciphertext)
	nonceHex := hex.EncodeToString(nonce)

	// Decode ciphertext and nonce from hexadecimal strings
	decodedCiphertext, err := hex.DecodeString(ciphertextHex)
	if err != nil {
		t.Errorf("Ciphertext decoding error: %v", err)
	}
	decodedNonce, err := hex.DecodeString(nonceHex)
	if err != nil {
		t.Errorf("Nonce decoding error: %v", err)
	}

	// Decrypt the decoded ciphertext
	decrypted, err := DecryptBytesGCM(decodedCiphertext, decodedNonce)
	if err != nil {
		t.Errorf("Decryption error: %v", err)
	}

	// Check if the decrypted plaintext matches the original plaintext
	if !bytes.Equal(decrypted, plaintext) {
		t.Errorf("Decrypted plaintext does not match the original plaintext")
	}
}

func TestEncryptDecryptFile(t *testing.T) {
	// Read the file content
	plaintext := []byte("This is a test")
	Init("/tmp")
	// Encrypt and write the content to a file
	err := writeFile(dbFile, plaintext)
	if err != nil {
		t.Fatal(err)
	}
	generateSecretKey()
	// Read and decrypt the content from the file
	decryptedContent, err := readFile(dbFile)
	if err != nil {
		t.Fatal(err)
	}

	// Compare the original plaintext and the decrypted content
	if !bytes.Equal(plaintext, decryptedContent) {
		t.Fatal(errors.New("Decryption failed: Plaintext and decrypted content do not match"))
	}
	os.Remove("/tmp/shift.db")
}

func TestWriteReadAccount(t *testing.T) {
	Init("/tmp")
	account = Account{Name: "Hans", Language: "en", Scooping: time.Date(2023, 12, 23, 9, 0, 0, 0, time.Local)}
	trans1 := Transaction{Amount: 13, Date: time.Date(2023, 3, 4, 0, 0, 0, 0, time.Local), Typ: Lmp}
	trans2 := Transaction{Amount: 15, Date: time.Date(2021, 4, 2, 0, 0, 0, 0, time.Local), Typ: Lmp}
	account.Transactions = append(account.Transactions, trans1, trans2)
	WriteAccount()

	tempAccount := Account{Name: "Hans", Language: "en", Scooping: time.Date(2023, 12, 23, 9, 0, 0, 0, time.Local)}
	tempAccount.Transactions = append(tempAccount.Transactions, trans1, trans2)

	account = Account{}
	if ReadAccount() != true {
		t.Error("Account could not be read")
	}
	if !reflect.DeepEqual(tempAccount, account) {
		t.Errorf("Account mismatch:\nExpected: %v\nGot: %v", tempAccount, account)
	}
	os.Remove("/tmp/shift.db")
}

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
