package lib

import (
	"bytes"
	"encoding/gob"
	"log"
)

var peerList []_peer

type _peer struct {
	Name           string
	CryptoKey      []byte // first peer in the list is local and a private key, for all others its a public key
	StorjBucket    string
	StorjAccessKey string
}

func addPeer(name string, publicKey []byte, storjBucket string, storjAccessKey string) {
	peer := _peer{Name: name, CryptoKey: publicKey, StorjBucket: storjBucket, StorjAccessKey: storjAccessKey}
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
