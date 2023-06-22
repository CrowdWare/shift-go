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
	peer := _peer{Uuid: "1234", CryptoKey: []byte("pubkey"), StorjBucket: "bucket", StorjAccessToken: "acckey"}
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
	peer := _peer{Uuid: "1234", CryptoKey: []byte("pubkey"), StorjBucket: "bucket", StorjAccessToken: "acckey"}
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

	if peerList[0].Uuid != "1234" {
		t.Errorf("Expected peer 1234 to be Hans but got %s", peerList[0].Uuid)
	}

	os.Remove("/tmp/peers.db")
}

func TestAddPeer(t *testing.T) {
	peerFile = "/tmp/peers.db"
	peerList = []_peer{}
	peer := _peer{Uuid: "1234", CryptoKey: []byte("pubkey"), StorjBucket: "bucket", StorjAccessToken: "acckey"}
	peerList = append(peerList, peer)
	peer2 := _peer{Uuid: "1235", CryptoKey: []byte("pubkey"), StorjBucket: "bucket", StorjAccessToken: "acckey"}
	peerList = append(peerList, peer2)
	writePeers()

	addPeer("Hans", "1235", []byte("pubkeyNew"), "newbucket", "newtoken")

	if len(peerList) != 2 {
		t.Errorf("Expected len to be 2 but got %d", len(peerList))
	}

	if peerList[1].StorjBucket != "newbucket" {
		t.Errorf("Expected newbucket but found %s", peerList[2].StorjBucket)
	}
}
