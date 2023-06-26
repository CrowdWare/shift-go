package lib

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"log"
	"time"

	"storj.io/uplink"
)

var messageMap map[string]_message

type _message struct {
	From     string
	Message  string
	PeerUuid string
	Time     time.Time
}

func createMessages() {
	messageMap = make(map[string]_message)
	writeMessages()
}

func addMessage(key, from, message, peerUuid string, time time.Time) {
	if existingMessage, ok := messageMap[key]; ok {
		// Key already exists, update the existing message
		existingMessage.From = from
		existingMessage.Message = message
		existingMessage.PeerUuid = peerUuid
		existingMessage.Time = time
		messageMap[key] = existingMessage
	} else {
		// Key doesn't exist, add a new message
		messageMap[key] = _message{
			From:     from,
			Message:  message,
			PeerUuid: peerUuid,
			Time:     time,
		}
	}
	writeMessages()
}

func readMessages() bool {
	buffer, err := readFile(messageFile)
	if err != nil {
		if debug {
			log.Println("Error reading file " + messageFile + ", " + err.Error())
		}
		return false
	}
	decoder := gob.NewDecoder(bytes.NewReader(buffer))
	err = decoder.Decode(&messageMap)
	if err != nil {
		if debug {
			log.Println("readPeers: " + err.Error())
		}
		return false
	}
	return true
}

func writeMessages() {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(messageMap)
	if err != nil {
		if debug {
			log.Println(err)
		}
		return
	}
	err = writeFile(messageFile, buffer.Bytes())
	if err != nil {
		if debug {
			log.Println(err)
		}
		return
	}
}

func getMessagesfromPeer(peerUuid string) ([]string, error) {
	var emptyList = make([]string, 0)
	ctx := context.Background()

	access, err := uplink.ParseAccess(peerMap[account.Uuid].StorjAccessToken)
	if err != nil {
		if debug {
			log.Printf("parse access failed %s", err.Error())
		}
		return emptyList, err
	}
	keys, err := listObjects(peerMap[account.Uuid].StorjBucket, "shift/messages/"+peerUuid+"/", ctx, access)
	if err != nil {
		if debug {
			log.Printf("list oebjects failed %s", err.Error())
		}
		return emptyList, err
	}
	return keys, nil
}

func getPeerMessage(peerUuid, key string) (string, time.Time, error) {
	ctx := context.Background()
	var t time.Time

	access, err := uplink.ParseAccess(peerMap[account.Uuid].StorjAccessToken)
	if err != nil {
		if debug {
			log.Printf("parse access failed %s", err.Error())
		}
		return "", t, err
	}
	ciphertext, time, err := get(key, peerMap[account.Uuid].StorjBucket, ctx, access)
	if err != nil {
		if debug {
			log.Printf("get failed %s", err.Error())
		}
		return "", t, err
	}

	plaintext, err := decryptString(peerMap[account.Uuid].CryptoKey, ciphertext)
	if err != nil {
		if debug {
			log.Printf("decrypt failed %s", err.Error())
		}
		return "", t, err
	}
	return plaintext, time, nil
}

func doesPeerMessageExist(peerUuid, messageKey string) (bool, error) {
	peer, ok := peerMap[peerUuid]
	if !ok {
		return false, fmt.Errorf("Peer not found in map")
	}
	ctx := context.Background()

	access, err := uplink.ParseAccess(peer.StorjAccessToken)
	if err != nil {
		return false, err
	}
	exists, err := exists(messageKey, peer.StorjBucket, ctx, access)
	if err != nil {
		return false, err
	}
	if exists {
		return true, nil
	}
	return false, nil
}

func deletePeerMassage(peerUuid, messageKey string) (bool, error) {
	peer, ok := peerMap[peerUuid]
	if !ok {
		return false, fmt.Errorf("Peer does not exist")
	}
	ctx := context.Background()

	access, err := uplink.ParseAccess(peer.StorjAccessToken)
	if err != nil {
		return false, err
	}
	err = delete(messageKey, peer.StorjBucket, ctx, access)
	if err != nil {
		return false, err
	}
	return true, nil
}
