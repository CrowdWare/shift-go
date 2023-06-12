package lib

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"log"
	"os"
	"time"
)

func generateSecretKey(db bool) ([]byte, error) {
	// Declare the variables used for key derivation
	var value3 int
	var empty []byte

	// Check if the file exists
	if db {
		if fileInfo, err := os.Stat(dbFile); os.IsNotExist(err) {
			// Use the current date as a fallback value
			value3 = time.Now().Day()
		} else if err != nil {
			return empty, err
		} else {
			// Retrieve the file's modification time
			modTime := fileInfo.ModTime()
			value3 = modTime.Day()
		}
	}
	// Define the variables used for key derivation
	var variable1 = 7539
	var variable2 = 2375

	// Perform calculations to derive the key
	derivedKey := int64(variable1 * variable2 * value3 * 9736635 * 42367736)

	// Convert the derived key to a 32-byte slice
	key := make([]byte, 32)
	for i := 0; i < 8; i++ {
		key[i] = byte(derivedKey >> (i * 8))
	}
	derivedKey = int64(variable1 + variable2*value3*9736635*42367736)
	for i := 8; i < 16; i++ {
		key[i] = byte(derivedKey >> ((i - 8) * 8))
	}
	derivedKey = int64(variable1 + variable2 + value3*9736635*42367736)
	for i := 16; i < 24; i++ {
		key[i] = byte(derivedKey >> ((i - 16) * 8))
	}
	derivedKey = int64(variable1*variable2 + value3*9736615*42367736)
	for i := 24; i < 32; i++ {
		key[i] = byte(derivedKey >> ((i - 24) * 8))
	}

	return key, nil
}

func EncryptStringGCM(value string) string {
	key, err := generateSecretKey(false)
	if err != nil {
		log.Fatal(err)
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatal(err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		log.Fatal(err)
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		log.Fatal(err)
	}

	ciphertext := aesGCM.Seal(nonce, nonce, []byte(value), nil)

	return hex.EncodeToString(ciphertext)
}

func EncryptBytesGCM(plaintext []byte) ([]byte, []byte, error) {
	key, err := generateSecretKey(true)
	if err != nil {
		log.Fatal(err)
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, nil, err
	}

	ciphertext := aesGCM.Seal(nil, nonce, plaintext, nil)

	return ciphertext, nonce, nil
}

func DecryptStringGCM(value string) string {
	key, err := generateSecretKey(false)
	if err != nil {
		log.Fatal(err)
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatal(err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		log.Fatal(err)
	}

	encryptedData, err := hex.DecodeString(value)
	if err != nil {
		log.Fatal(err)
	}

	iv := encryptedData[:12] // GCM IV is usually 12 bytes
	cipherText := encryptedData[12:]

	plaintext, err := aesGCM.Open(nil, iv, cipherText, nil)
	if err != nil {
		log.Fatal(err)
	}
	return string(plaintext)
}

func DecryptBytesGCM(ciphertext, nonce []byte) ([]byte, error) {
	key, err := generateSecretKey(true)
	if err != nil {
		log.Fatal(err)
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
