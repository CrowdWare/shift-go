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
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/gob"
	"encoding/pem"
	"fmt"
	"log"
)

var peerMap map[string]_peer

type _peer struct {
	Uuid             string
	Name             string
	CryptoKey        []byte // first peer in the list is local and a private key, for all others its a public key
	StorjBucket      string
	StorjAccessToken string
}

func createPeer() int {
	peerMap = make(map[string]_peer)
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		if debug {
			fmt.Println("Failed to generate RSA key pair:", err)
		}
		return 1
	}
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})
	log.Println("Create peer: " + account.Name + ", " + account.Uuid)
	localPeer := _peer{Uuid: account.Uuid, Name: account.Name, CryptoKey: privateKeyPEM, StorjBucket: "", StorjAccessToken: ""}
	peerMap[account.Uuid] = localPeer
	writePeers()
	return 0
}

func addPeer(name string, uuid string, publicKey []byte, storjBucket string, storjAccessToken string) {
	existingPeer, ok := peerMap[uuid]
	if ok {
		// update peer
		existingPeer.CryptoKey = publicKey
		existingPeer.StorjBucket = storjBucket
		existingPeer.StorjAccessToken = storjAccessToken
		peerMap[uuid] = existingPeer
	} else {
		// append peer
		peerMap[uuid] = _peer{Uuid: uuid, Name: name, CryptoKey: publicKey, StorjBucket: storjBucket, StorjAccessToken: storjAccessToken}
	}
	writePeers()
}

func readPeers() bool {
	buffer, err := readFile(peerFile)
	if err != nil {
		if debug {
			log.Println("Error reading file " + peerFile + ", " + err.Error())
		}
		return false
	}
	decoder := gob.NewDecoder(bytes.NewReader(buffer))
	err = decoder.Decode(&peerMap)
	if err != nil {
		if debug {
			log.Println("readPeers: " + err.Error())
		}
		return false
	}
	return true
}

func writePeers() {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(peerMap)
	if err != nil {
		if debug {
			log.Println(err)
		}
		return
	}
	err = writeFile(peerFile, buffer.Bytes())
	if err != nil {
		if debug {
			log.Println(err)
		}
		return
	}
}
