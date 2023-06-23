package lib

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"log"
	"os"
	"time"
)

func generateSecretKey(db bool, read bool, fileName string) ([]byte, error) {
	var variable1 = int64(var1)
	var variable2 = int64(var2)
	var variable3 = int64(var3)
	var variable4 = int64(var4)
	var variable5 = int64(var5)

	// Check if the file exists
	if db {
		variable3 = int64(time.Now().Day())
		fileInfo, err := os.Stat(fileName)
		if read && err == nil {
			modTime := fileInfo.ModTime()
			variable3 = int64(modTime.Day())
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
		key, err = generateSecretKey(false, false, "")
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

func encryptBytesGCM(plaintext []byte, fileName string) ([]byte, []byte, error) {
	key, err := generateSecretKey(true, false, fileName)
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
		key, err = generateSecretKey(false, false, "")
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

func decryptBytesGCM(ciphertext, nonce []byte, fileName string) ([]byte, error) {
	key, err := generateSecretKey(true, true, fileName)
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

func encryptString(publicKeyBytes []byte, plainText string) ([]byte, error) {
	// Retrieve the public key from bytes
	parsedPublicKey, err := x509.ParsePKIXPublicKey(publicKeyBytes)
	if err != nil {
		if debug {
			log.Println("Error parsing PKIXPub, " + err.Error())
		}
		return nil, err
	}

	// Convert the parsed public key to the correct type
	retrievedPublicKey := parsedPublicKey.(*rsa.PublicKey)

	// Encrypt the plaintext using the public key
	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, retrievedPublicKey, []byte(plainText))
	if err != nil {
		if debug {
			log.Println("Error encryptPKCS" + err.Error())
		}
		return nil, err
	}
	return ciphertext, nil
}

func decryptString(privateKeyBytes []byte, ciphertext []byte) (string, error) {
	// Parse the private key from bytes
	block, _ := pem.Decode(privateKeyBytes)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return "", fmt.Errorf("Failed to decode PEM block containing private key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Println(err)
		return "", err
	}

	// Decrypt the ciphertext using the private key
	decryptedText, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, ciphertext)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return string(decryptedText), nil
}
