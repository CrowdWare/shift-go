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

	localPeer := _peer{Uuid: account.Uuid, Name: account.Name, CryptoKey: privateKeyPEM, StorjBucket: "", StorjAccessToken: ""}
	peerMap[account.Uuid] = localPeer
	writePeers()
	return 0
}

func addPeer(name string, uuid string, publicKey []byte, storjBucket string, storjAccessToken string) {
	if existingPeer, ok := peerMap[uuid]; ok {
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
