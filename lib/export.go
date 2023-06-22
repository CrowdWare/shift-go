package lib

import (
	"context"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"flag"
	"fmt"
	"log"
	"runtime"
	"strings"
	"time"

	"github.com/google/uuid"
	"storj.io/uplink"
)

/*
**	This is needed to create a json string to be sent to the client.
 */
type TransactionTO struct {
	Pkey    string
	Amount  int64
	Date    int64
	From    string
	To      string
	Purpose string
	Typ     int
}

/*
**	Set the path where the db can be stored.
 */
func Init(filesDir string) {
	// avoid scooping on a desktop
	if runtime.GOOS != "android" && runtime.GOOS != "ios" {
		growLevel0 = 0
		growLevel1 = 0
		growLevel2 = 0
		growLevel3 = 0
	}
	dbFile = filesDir + "/shift.db"
	peerFile = filesDir + "/peers.db"
	if !readPeers() {
		createPeer()
	}
}

/*
**	Return the account Uuid
 */
func GetUuid() string {
	return account.Uuid
}

/*
**	Return the account Name
 */
func GetName() string {
	return account.Name
}

/*
**	Set the account Name
 */
func SetName(name string) {
	account.Name = name
	writeAccount()
}

/*
**	Checks if account has been created already.
 */
func HasJoined() bool {
	return readAccount()
}

/*
**	Checks if account is scooping and is already 20 hours ago.
**	If so a new scooped transaction is added and isScooping is set to false.
**	The function returns this isScooping flag.
 */
func IsScooping() bool {
	return checkScooping()
}

/*
**	Create the account and send it to the rest service.
** 	This is only called once running the app for the first time.
 */
func CreateAccount(name, uuid, ruuid, country, language string) int {
	res := addAccount(name, uuid, ruuid, country, language, false)
	cp := createPeer()
	return res + cp
}

/*
**	Get a list of all mates from rest service, pack them into json and return the json string.
**	We are also adding the user from the peerlist. These are the users with whom we exchanged
**	Our public key and storj access with.
 */
func GetMatelist() string {
	list := getMatelist(false)
	for _, p := range peerList {
		index := contains(list, p.Uuid)
		if index > 0 {
			list[index].HasPeerData = true
		} else {
			list = append(list, Friend{Name: p.Name, Uuid: p.Uuid, HasPeerData: true})
		}
	}
	jsonData, err := json.Marshal(list)
	if err != nil {
		return ""
	}
	return string(jsonData)
}

/*
**	Set the account in scooping mode, scooping is set to now and rest method is called.
**	The level counts will be set after the rest call.
 */
func StartScooping() int {
	if account.IsScooping {
		return 1
	}
	return setScooping(false)
}

/*
**	Sum up all transactions after subtracting the demurrage.
 */
func GetBalanceInMillis() int64 {
	balance := int64(0)
	for _, t := range account.Transactions {
		if t.Typ == Lmp {
			// LMP is a withdraw, so its a negative value
			// the demurrage is calculated on the receivers side
			balance += t.Amount * 1000
		} else {
			balance += calculateWorthInMillis(t.Amount, t.Date)
		}
	}
	return balance
}

/*
**	Calculate the scooped amount which since last setScooping, to be displayed
**	on the client while balance display is in milli liter mode.
 */
func GetScoopedBalance() int64 {
	if account.IsScooping {
		diff := time.Now().Sub(account.Scooping)
		return calcGrowPerDiff(diff)
	}
	return 0
}

/*
**	Returns the amount of hours the account is scooping
 */
func GetScoopingHours() float64 {
	if account.IsScooping {
		diff := time.Now().Sub(account.Scooping)
		return diff.Hours()
	}
	return 0
}

/*
**	Get a list of the last 30 transactions, pack it into a json string and return it.
 */
func GetTransactions() string {
	trans := make([]TransactionTO, 0)
	startIndex := len(account.Transactions) - 30
	if startIndex < 0 {
		startIndex = 0
	}
	for i := len(account.Transactions) - 1; i >= startIndex; i-- {
		t := account.Transactions[i]
		trans = append(trans, TransactionTO{Pkey: t.Pkey, Amount: t.Amount, Purpose: t.Purpose, Date: t.Date.Unix(), From: t.From, To: t.To, Typ: int(t.Typ)})
	}
	jsonData, err := json.Marshal(trans)
	if err != nil {
		return ""
	}
	return string(jsonData)
}

/*
**	When GetProposalFromQRCode is called, then this transaction is stored temporary
**  In this method we take the last transaction and withdraw it from the account.
**  We don't want to have a public function for adding transactions.
 */
func AcceptProposal() string {
	lastTransaction.Pkey = uuid.New().String()
	err := addTransaction(lastTransaction.Pkey, lastTransaction.Amount, lastTransaction.Purpose, lastTransaction.Date, lastTransaction.From, lastTransaction.To, lastTransaction.Typ, lastTransaction.Uuid)
	if err != nil {
		return err.Error()
	}
	if lastTransaction.Purpose == "" {
		if debug {
			log.Println("AcceptProposal no purpose")
		}
		return "NO_PURPOSE"
	}
	return "ok"
}

/*
**	Create a transaction, pack it to a json string and return the encrypted string.
**  On the client a QR code will be created based on that string.
 */
func GetProposalQRCode(amount int64, purpose string) string {
	trans := _transaction{Pkey: uuid.New().String(), Amount: amount * -1, Purpose: purpose, Date: time.Now(), Typ: Lmp, From: account.Name, To: account.Name, Uuid: account.Uuid}
	jsonData, err := json.Marshal(trans)
	if err != nil {
		if debug {
			log.Println(err)
		}
		return ""
	}
	if purpose == "" {
		if debug {
			log.Println("GetProposalQRCode no purpose")
		}
		return "NO_PURPOSE"
	}
	return encryptStringGCM(string(jsonData), false)
}

/*
**	Create a transation, pack it to a json string and return the encrypted string.
**	On the client a QR code will be created based on that string.
 */
func GetAgreementQRCode() string {
	trans := _transaction{Pkey: lastTransaction.Pkey, Amount: lastTransaction.Amount * -1, Purpose: lastTransaction.Purpose, Date: lastTransaction.Date, Typ: Lmr, From: account.Name, To: lastTransaction.To, Uuid: lastTransaction.Uuid}
	jsonData, err := json.Marshal(trans)
	if err != nil {
		if debug {
			log.Println(err)
		}
		return ""
	}
	if lastTransaction.Purpose == "" {
		if debug {
			log.Println("GetAgreementQRCode no purpose")
		}
		return "NO_PURPOSE"
	}
	return encryptStringGCM(string(jsonData), false)
}

/*
**	Get a transaction from account, pack it to a json string and return the encrypted string.
**	This is used in case the giver skipped the last step and wants to recreate the qr code.
 */
func GetAgreementQRCodeForTransaction(pkey string) string {
	found := false
	isLmp := false
	var trans _transaction
	for _, t := range account.Transactions {
		if t.Pkey == pkey {
			trans.Pkey = t.Pkey
			trans.Date = t.Date
			if t.Typ == Lmp {
				isLmp = true
				trans.Amount = t.Amount * -1
				trans.From = account.Name
				trans.Typ = Lmr
			} else {
				isLmp = false
				trans.Amount = t.Amount
				trans.From = t.From
				trans.Typ = t.Typ
			}
			trans.To = t.To
			trans.Purpose = t.Purpose
			trans.Uuid = t.Uuid
			found = true
			break
		}
	}
	if found == false {
		return "|NOT FOUND"
	}
	jsonData, err := json.Marshal(trans)
	if err != nil {
		if debug {
			log.Println(err)
		}
		return " | "
	}
	if isLmp {
		return string(jsonData) + "|" + encryptStringGCM(string(jsonData), false)
	}
	return string(jsonData) + "|NOT LMP"
}

/*
**	We decrypt the string from a QR code, unpack and validate the transaction
**	and send back the json to the client (we cannot send objects to Kotlin).
**  If it cannot be decrypted or unpacked its a sign for an attack.
 */
func GetProposalFromQRCode(enc string) string {
	jsonData, err := decryptStringGCM(enc, false)
	if err != nil {
		if debug {
			log.Println("BookTransaction: error decrypting transaction")
		}
		return "FRAUD"
	}

	err = json.Unmarshal([]byte(jsonData), &lastTransaction)
	if err != nil {
		if debug {
			log.Println("BookTransaction: error unmarshaling transaction")
		}
		return "FRAUD"
	}
	if lastTransaction.Typ != Lmp {
		return "WRONG_TYP"
	}
	if lastTransaction.Typ == Lmp && lastTransaction.Amount > 0 {
		if debug {
			log.Println("BookTransaction: amount > 0")
		}
		return "FRAUD"
	}
	if lastTransaction.Purpose == "" {
		if debug {
			log.Println("GetProposalFromQRCode no purpose")
		}
		return "NO_PURPOSE"
	}

	return string(jsonData)
}

/*
**	We encrypt the QR-Code, unpack and validate the transaction and when
**	everything is fine we book it.
 */
func GetAgreementFromQRCode(enc string) string {
	jsonData, err := decryptStringGCM(enc, false)
	if err != nil {
		if debug {
			log.Println("GetAgreementFromQRCode: error decrypting transaction")
		}
		return "FRAUD"
	}
	err = json.Unmarshal([]byte(jsonData), &lastTransaction)
	if err != nil {
		if debug {
			log.Println("GetAgreementFromQRCode: error unmarshaling transaction")
		}
		return "FRAUD"
	}
	if lastTransaction.Typ == Lmp {
		return "WRONG_TYP"
	}
	if lastTransaction.Typ != Lmr {
		if debug {
			log.Println("GetAgreementFromQRCode: wrong transaction typ")
		}
		return "FRAUD"
	}
	if lastTransaction.Amount < 1 {
		if debug {
			log.Println("GetAgreementFromQRCode: error amount < 1")
		}
		return "WRONG_TYP"
	}
	if lastTransaction.Uuid != account.Uuid {
		if debug {
			log.Println("GetAgreementFromQRCode: error transaction not for this user")
		}
		return "WRONG_RECEIVER"
	}
	if transactionExists(lastTransaction.Pkey) {
		if debug {
			log.Println("GetAgreementFromQRCode: transaction already booked")
		}
		return "DOUBLE_SPENT"
	}
	if lastTransaction.Purpose == "" {
		if debug {
			log.Println("GetAgreementFromQRCode no purpose")
		}
		return "NO_PURPOSE"
	}
	// only check error on withdraw, not on receive
	addTransaction(lastTransaction.Pkey, lastTransaction.Amount, lastTransaction.Purpose, lastTransaction.Date, lastTransaction.From, lastTransaction.To, lastTransaction.Typ, lastTransaction.Uuid)

	return "ok"
}

/*
**	Decrypt the QR-Code, unmarshalls the peer and add it to the peer list
 */
func AddPeerFromQRCode(enc string) bool {
	jsonData, err := decryptStringGCM(enc, false)
	if err != nil {
		if debug {
			log.Println(err)
		}
		return false
	}
	var peer _peer
	err = json.Unmarshal([]byte(jsonData), &peer)
	if err != nil {
		if debug {
			log.Println("GetAgreementFromQRCode: error unmarshaling transaction")
		}
		return false
	}
	addPeer(peer.Name, peer.Uuid, peer.CryptoKey, peer.StorjBucket, peer.StorjAccessToken)
	return true
}

/*
**	Returns an encrypted record of the first Peer
 */
func GetPeerQRCode() string {
	if len(peerList) == 0 {
		cp := createPeer()
		if cp != 0 {
			return ""
		}
	}
	if peerList[0].StorjBucket == "" {
		return ""
	}
	// Decode the private key from the string
	block, _ := pem.Decode(peerList[0].CryptoKey)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		if debug {
			fmt.Println("Failed to decode private key")
		}
		return ""
	}
	// Parse the private key
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		if debug {
			fmt.Println("Failed to parse private key:", err)
		}
		return ""
	}

	// Get the public key from the private key
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		fmt.Println("Failed to encode public key:", err)
		return ""
	}

	peer := _peer{Name: account.Name, Uuid: account.Uuid, CryptoKey: publicKeyBytes, StorjBucket: peerList[0].StorjBucket, StorjAccessToken: peerList[0].StorjAccessToken}
	jsonData, err := json.Marshal(peer)
	if err != nil {
		if debug {
			log.Println(err)
		}
		return ""
	}
	return encryptStringGCM(string(jsonData), false)
}

/*
**	Saves the Storj Data
 */
func SetStorj(bucketName string, accessKey string) bool {
	if bucketName == "" || accessKey == "" {
		return false
	}
	peerList[0].StorjBucket = bucketName
	peerList[0].StorjAccessToken = accessKey
	writePeers()
	return true
}

/*
**	Returns the Storj Bucketname
 */
func GetBucketName() string {
	if len(peerList) > 0 {
		return peerList[0].StorjBucket
	}
	return ""
}

/*
**	Return the Storj Access Token
 */
func GetAccessToken() string {
	if len(peerList) > 0 {
		return peerList[0].StorjAccessToken
	}
	return ""
}

/*
**	Puts an encrypted message on the Storj bucket from the peer
 */
func SendMessageToPeer(peerUuid string, message string) string {
	peer := -1
	for i, p := range peerList {
		if p.Uuid == peerUuid {
			peer = i
		}
	}
	if peer == -1 {
		return "1"
	}
	ctx := context.Background()
	accessGrant := flag.String("accessSend", peerList[peer].StorjAccessToken, "access grant from satellite")
	access, err := uplink.ParseAccess(*accessGrant)
	if err != nil {
		if debug {
			log.Printf("parse access failed %s", err.Error())
		}
		return "2"
	}

	cipherText, err := encryptString(peerList[peer].CryptoKey, message)
	if err != nil {
		if debug {
			log.Println("error encryting the message: " + err.Error())
		}
		return "3"
	}

	messageKey := "shift/messages/" + peerList[0].Uuid + "/" + uuid.NewString()
	err = put(messageKey, cipherText, peerList[peer].StorjBucket, ctx, access)
	if err != nil {
		if debug {
			log.Println("put failed: " + err.Error())
		}
		return "4"
	}
	return messageKey
}

/*
**	Returns a comma separated list of keys. These keys can be used to "get" a message from Storj
 */
func GetMessagesfromPeer(peerUuid string) string {
	ctx := context.Background()
	accessGrant := flag.String("accessList", peerList[0].StorjAccessToken, "access grant from satellite")
	access, err := uplink.ParseAccess(*accessGrant)
	if err != nil {
		if debug {
			log.Printf("parse access failed %s", err.Error())
		}
		return "1"
	}
	keys, err := listObjects(peerList[0].StorjBucket, "shift/messages/"+peerUuid+"/", ctx, access)
	if err != nil {
		if debug {
			log.Printf("list oebjects failed %s", err.Error())
		}
		return "2"
	}

	return strings.Join(keys, ",")
}

/*
**	Gets the message from Storj and returns it decrypted
 */
func GetPeerMessage(peerUuid, key string) string {
	ctx := context.Background()
	accessGrant := flag.String("accessGet", peerList[0].StorjAccessToken, "access grant from satellite")
	access, err := uplink.ParseAccess(*accessGrant)
	if err != nil {
		if debug {
			log.Printf("parse access failed %s", err.Error())
		}
		return "1"
	}
	ciphertext, err := get(key, peerList[0].StorjBucket, ctx, access)
	if err != nil {
		if debug {
			log.Printf("get failed %s", err.Error())
		}
		return "2"
	}

	plaintext, err := decryptString(peerList[0].CryptoKey, ciphertext)
	if err != nil {
		if debug {
			log.Printf("decrypt failed %s", err.Error())
		}
		return "3"
	}
	return "0," + plaintext
}
