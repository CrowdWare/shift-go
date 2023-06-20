package lib

import (
	"bytes"
	"encoding/gob"
	"os"
	"reflect"
	"testing"
)

func TestPeerSerialize(t *testing.T) {
	var peerlist []_peer
	peer := _peer{Name: "Hans", CryptoKey: []byte("pubkey"), StorjBucket: "bucket", StorjAccessKey: "acckey"}
	peerlist = append(peerlist, peer)
	var buffer bytes.Buffer

	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(peerlist)
	if err != nil {
		t.Error(err)
	}
	var peerlist2 []_peer
	decoder := gob.NewDecoder(bytes.NewReader(buffer.Bytes()))
	err = decoder.Decode(&peerlist2)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(peerlist, peerlist2) {
		t.Errorf("Peer mismatch:\nExpected: %v\nGot: %v", peerlist, peerlist2)
	}
}

func TestReadWritePeers(t *testing.T) {
	peerFile = "/tmp/peers.db"
	peerList = []_peer{}
	peer := _peer{Name: "Hans", CryptoKey: []byte("pubkey"), StorjBucket: "bucket", StorjAccessKey: "acckey"}
	peerList = append(peerList, peer)
	writePeers()

	result := readPeers()
	expected := true
	if result != expected {
		t.Errorf("Unexpected result. Got: %v, Expected: %v", result, expected)
	}
	account = _account{}
	writePeers()

	result = readPeers()
	expected = true
	if result != expected {
		t.Errorf("Unexpected result. Got: %v, Expected: %v", result, expected)
	}

	if len(peerList) != 1 {
		t.Errorf("Expected peercount to be 1 but got %d", len(peerList))
	}

	if peerList[0].Name != "Hans" {
		t.Errorf("Expected peer name to be Hans but got %s", peerList[0].Name)
	}

	os.Remove("/tmp/peers.db")
}
