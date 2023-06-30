/****************************************************************************
 * Copyright (C) 2023 CrowdWare
 *
 * This file is part of SHIFT.
 *
 *  SHIFT is free software: you can redistribute it and/or modify
 *  it under the terms of the GNU General Public License as published by
 *  the Free Software Foundation, either version 3 of the License, or
 *  (at your option) any later version.
 *
 *  SHIFT is distributed in the hope that it will be useful,
 *  but WITHOUT ANY WARRANTY without even the implied warranty of
 *  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *  GNU General Public License for more details.
 *
 *  You should have received a copy of the GNU General Public License
 *  along with SHIFT.  If not, see <http://www.gnu.org/licenses/>.
 *
 ****************************************************************************/
package lib

import (
	"bytes"
	"encoding/hex"
	"errors"
	"os"
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	dbFile = "/tmp/shift.db"
	teststring := "The quick brown fox"
	enc := encryptStringGCM(teststring, false)
	result, err := decryptStringGCM(enc, false)
	if err != nil {
		t.Error(err)
		return
	}
	expected := teststring
	if result != expected || enc == teststring {
		t.Errorf("Unexpected result. Got: %s, Expected: %s", result, expected)
	}

	result, err = decryptStringGCM(enc+"a", false)
	if err == nil {
		t.Error("Expected an error decrypting")
	}
	runes := []rune(enc)
	runes[3] = 'a'
	result, err = decryptStringGCM(string(runes), false)
	if err == nil {
		t.Error("Expected an error decrypting")
	}
}

func TestEncryptAndDecryptBytesGCM(t *testing.T) {
	plaintext := []byte("Hello, World!")
	dbFile = "/tmp/shift.db"
	// Encrypt the plaintext
	ciphertext, nonce, err := encryptBytesGCM(plaintext, dbFile)
	if err != nil {
		t.Errorf("Encryption error: %v", err)
	}

	// Decrypt the ciphertext
	decrypted, err := decryptBytesGCM(ciphertext, nonce, dbFile)
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
	dbFile = "/tmp/shift.db"
	// Encrypt the plaintext
	ciphertext, nonce, err := encryptBytesGCM(plaintext, dbFile)
	if err != nil {
		t.Errorf("Encryption error: %v", err)
	}

	// Decrypt the ciphertext
	decrypted, err := decryptBytesGCM(ciphertext, nonce, dbFile)
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
	dbFile = "/tmp/shift.db"
	// Encrypt the plaintext
	ciphertext, nonce, err := encryptBytesGCM(plaintext, dbFile)
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
	decrypted, err := decryptBytesGCM(decodedCiphertext, decodedNonce, dbFile)
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
	dbFile = "/tmp/shift.db"
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
