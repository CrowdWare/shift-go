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
	"io/ioutil"
	"os"
	"sync"
)

var mutex = &sync.Mutex{} // Create a mutex lock

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
	mutex.Lock()
	defer mutex.Unlock()

	tempFilePath := filename + ".temp"
	ciphertext, nonce, err := encryptBytesGCM(content, filename)
	if err != nil {
		return err
	}

	// Append the nonce to the ciphertext
	contentWithNonce := append(nonce, ciphertext...)

	// Write the encrypted content with nonce to a file
	if err := ioutil.WriteFile(tempFilePath, contentWithNonce, 0644); err != nil {
		return err
	}

	// Sync the file to ensure data is written to disk
	file, err := os.OpenFile(tempFilePath, os.O_WRONLY, 0)
	if err != nil {
		return err
	}
	defer file.Close()

	err = file.Sync()
	if err != nil {
		return err
	}

	// Replace the original file with the decrypted file
	err = os.Rename(tempFilePath, filename)
	if err != nil {
		return err
	}

	return nil
}

func readFile(filename string) ([]byte, error) {
	mutex.Lock()
	defer mutex.Unlock()

	contentWithNonce, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// Extract the nonce and ciphertext from the content
	nonce := contentWithNonce[:12]
	ciphertext := contentWithNonce[12:]

	plainbytes, err := decryptBytesGCM(ciphertext, nonce, filename)
	if err != nil {
		return nil, err
	}

	return plainbytes, nil
}
