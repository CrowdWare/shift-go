package lib

import (
	"bytes"
	"encoding/gob"
	"os"
	"reflect"
	"testing"
)

func TestPeerSerialize(t *testing.T) {
	peerMap = map[string]_peer{}
	peerMap["1234"] = _peer{Uuid: "1234", CryptoKey: []byte("pubkey"), StorjBucket: "bucket", StorjAccessToken: "acckey"}
	var buffer bytes.Buffer

	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(peerMap)
	if err != nil {
		t.Error(err)
	}
	var peerMap2 = map[string]_peer{}
	decoder := gob.NewDecoder(bytes.NewReader(buffer.Bytes()))
	err = decoder.Decode(&peerMap2)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(peerMap, peerMap2) {
		t.Errorf("Peer mismatch:\nExpected: %v\nGot: %v", peerMap, peerMap2)
	}
}

func TestReadWritePeers(t *testing.T) {
	peerFile = "/tmp/peers.db"
	peerMap = map[string]_peer{}
	peerMap["1234"] = _peer{Name: "Hans", CryptoKey: []byte("pubkey"), StorjBucket: "bucket", StorjAccessToken: "acckey"}
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

	if len(peerMap) != 1 {
		t.Errorf("Expected peercount to be 1 but got %d", len(peerMap))
	}

	if peer, ok := peerMap["1234"]; ok {
		if peer.Name != "Hans" {
			t.Errorf("Expected peer 1234 to be Hans but got %s", peer.Name)
		}
	}

	os.Remove("/tmp/peers.db")
}

func TestAddPeer(t *testing.T) {
	peerFile = "/tmp/peers.db"
	peerMap = map[string]_peer{}
	peerMap["1234"] = _peer{CryptoKey: []byte("pubkey"), StorjBucket: "bucket", StorjAccessToken: "acckey"}
	peerMap["1235"] = _peer{CryptoKey: []byte("pubkey"), StorjBucket: "bucket", StorjAccessToken: "acckey"}
	writePeers()

	addPeer("Hans", "1235", []byte("pubkeyNew"), "newbucket", "newtoken")

	if len(peerMap) != 2 {
		t.Errorf("Expected len to be 2 but got %d", len(peerMap))
	}

	if peerMap["1235"].StorjBucket != "newbucket" {
		t.Errorf("Expected newbucket but found %s", peerMap["1235"].StorjBucket)
	}
}
