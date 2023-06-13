package lib

import (
	"bytes"
	"encoding/hex"
	"errors"
	"os"
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	Init("")
	teststring := "The quick brown fox"
	enc := encryptStringGCM(teststring, false)
	result := decryptStringGCM(enc)
	expected := teststring
	if result != expected || enc == teststring {
		t.Errorf("Unexpected result. Got: %s, Expected: %s", result, expected)
	}
}

func TestEncryptAndDecryptBytesGCM(t *testing.T) {
	plaintext := []byte("Hello, World!")
	Init("")
	// Encrypt the plaintext
	ciphertext, nonce, err := encryptBytesGCM(plaintext)
	if err != nil {
		t.Errorf("Encryption error: %v", err)
	}

	// Decrypt the ciphertext
	decrypted, err := decryptBytesGCM(ciphertext, nonce)
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
	ciphertext, nonce, err := encryptBytesGCM(plaintext)
	if err != nil {
		t.Errorf("Encryption error: %v", err)
	}

	// Decrypt the ciphertext
	decrypted, err := decryptBytesGCM(ciphertext, nonce)
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
	ciphertext, nonce, err := encryptBytesGCM(plaintext)
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
	decrypted, err := decryptBytesGCM(decodedCiphertext, decodedNonce)
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
