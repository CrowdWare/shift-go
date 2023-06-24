package lib

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"
)

var messageMap map[string]_message

type _message struct {
	From     string
	Message  string
	PeerUuid string
	Time     time.Time
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
