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

var peerList []_peer

type _peer struct {
	Name             string
	Uuid             string
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

	localPeer := _peer{Name: account.Name, Uuid: account.Uuid, CryptoKey: privateKeyPEM, StorjBucket: "", StorjAccessToken: ""}
	peerList = append(peerList, localPeer)
	writePeers()
	return 0
}

func addPeer(name string, uuid string, publicKey []byte, storjBucket string, storjAccessToken string) {
	peer := _peer{Uuid: uuid, CryptoKey: publicKey, StorjBucket: storjBucket, StorjAccessToken: storjAccessToken}
	peerList = append(peerList, peer)
	writePeers()
}

func readPeers() bool {
	if fileExists(peerFile) {
		buffer, err := readFile(peerFile)
		if err != nil {
			if debug {
				log.Println("Error reading file " + peerFile + ", " + err.Error())
			}
			return true
		}
		decoder := gob.NewDecoder(bytes.NewReader(buffer))
		err = decoder.Decode(&peerList)
		if err != nil {
			if debug {
				log.Println("readPeers: " + err.Error())
			}
		}
		return true
	}
	return false
}

func writePeers() {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(peerList)
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
