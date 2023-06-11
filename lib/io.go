package lib

import (
	"io/ioutil"
	"os"
)

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	if err == nil {
		return true // File exists
	}
	if os.IsNotExist(err) {
		return false // File does not exist
	}
	return false // Error occurred while checking
}

func writeFile(filename string, content []byte) error {
	ciphertext, nonce, err := EncryptBytesGCM(content)
	if err != nil {
		return err
	}

	// Append the nonce to the ciphertext
	contentWithNonce := append(nonce, ciphertext...)

	// Write the encrypted content with nonce to a file
	if err := ioutil.WriteFile(filename, contentWithNonce, 0644); err != nil {
		return err
	}

	return nil
}

func readFile(filename string) ([]byte, error) {
	contentWithNonce, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// Extract the nonce and ciphertext from the content
	nonce := contentWithNonce[:12]
	ciphertext := contentWithNonce[12:]

	plaintext, err := DecryptBytesGCM(ciphertext, nonce)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
