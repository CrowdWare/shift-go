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

func generateSecretKey(db bool, read bool) ([]byte, error) {
	var variable1 = var1
	var variable2 = var2
	var variable3 = var3
	var variable4 = var4
	var variable5 = var5

	// Check if the file exists
	if db {
		variable3 = time.Now().Day()
		fileInfo, err := os.Stat(dbFile)
		if read && err == nil {
			modTime := fileInfo.ModTime()
			variable3 = modTime.Day()
		}
	}
	// Perform calculations to derive the key
	derivedKey := int64(variable1 * variable2 * variable3 * variable4 * variable5)

	// Convert the derived key to a 32-byte slice
	key := make([]byte, 32)
	for i := 0; i < 8; i++ {
		key[i] = byte(derivedKey >> (i * 8))
	}
	derivedKey = int64(variable1 + variable2*variable3*variable4*variable5)
	for i := 8; i < 16; i++ {
		key[i] = byte(derivedKey >> ((i - 8) * 8))
	}
	derivedKey = int64(variable1 + variable2 + variable3*variable4*variable5)
	for i := 16; i < 24; i++ {
		key[i] = byte(derivedKey >> ((i - 16) * 8))
	}
	derivedKey = int64(variable1*variable2 + variable3*variable4*variable5)
	for i := 24; i < 32; i++ {
		key[i] = byte(derivedKey >> ((i - 24) * 8))
	}

	return key, nil
}

func encryptStringGCM(value string, webservice bool) string {
	var key []byte
	var err error

	if webservice {
		var secret_key string
		secret_key, err = decryptStringGCM(secret_key_enc, false)
		key = []byte(secret_key)
	} else {
		key, err = generateSecretKey(false, false)
		if err != nil {
			if debug {
				log.Println(err)
			}
			return ""
		}
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		if debug {
			log.Println(err)
		}
		return ""
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		if debug {
			log.Println(err)
		}
		return ""
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		if debug {
			log.Println(err)
		}
		return ""
	}

	ciphertext := aesGCM.Seal(nonce, nonce, []byte(value), nil)

	return hex.EncodeToString(ciphertext)
}

func encryptBytesGCM(plaintext []byte) ([]byte, []byte, error) {
	key, err := generateSecretKey(true, false)
	if err != nil {
		if debug {
			log.Println("generateSecretKey" + err.Error())
		}
		return nil, nil, err
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

func decryptStringGCM(value string, webservice bool) (string, error) {
	var key []byte
	var err error

	if webservice {
		var secret_key string
		secret_key, err = decryptStringGCM(secret_key_enc, false)
		key = []byte(secret_key)
	} else {
		key, err = generateSecretKey(false, false)
		if err != nil {
			return "", err
		}
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	encryptedData, err := hex.DecodeString(value)
	if err != nil {
		return "", err
	}

	iv := encryptedData[:12] // GCM IV is usually 12 bytes
	cipherText := encryptedData[12:]

	plaintext, err := aesGCM.Open(nil, iv, cipherText, nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}

func decryptBytesGCM(ciphertext, nonce []byte) ([]byte, error) {
	key, err := generateSecretKey(true, true)
	if err != nil {
		if debug {
			log.Println(err)
		}
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	plainbytes, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plainbytes, nil
}
