package lib

import (
	"encoding/json"
	"log"
	"time"
)

/*
**	This is needed to create a json string to be sent to the client.
 */
type TransactionTO struct {
	Amount  int64
	Date    int64
	From    string
	Purpose string
	Typ     int
}

/*
**	Set the path where the db can be stored.
 */
func Init(filesDir string) {
	dbFile = filesDir + "/shift.db"
}

/*
**	Return the account Uuid
 */
func GetUuid() string {
	return account.Uuid
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
func CreateAccount(name, uuid, ruuid, country, language string) {
	addAccount(name, uuid, ruuid, country, language, false)
}

/*
**	Get a list of all mates from rest service, pack them into json and return the json string.
 */
func GetMatelist() string {
	list := getMatelist(false)
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
func StartScooping() {
	if account.IsScooping {
		return
	}
	account.IsScooping = true
	account.Scooping = time.Now()
	account.Level_1_count = 0
	account.Level_2_count = 0
	account.Level_3_count = 0
	writeAccount()
	setScooping(false)
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
		trans = append(trans, TransactionTO{Amount: t.Amount, Purpose: t.Purpose, Date: t.Date.Unix(), From: t.From, Typ: int(t.Typ)})
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
	err := addTransaction(lastTransaction.Amount, lastTransaction.Purpose, lastTransaction.Date, lastTransaction.From, lastTransaction.Typ, lastTransaction.Uuid)
	if err != nil {
		return err.Error()
	}
	if lastTransaction.Purpose == "" {
		log.Println("AcceptProposal no purpose")
		return "NO_PURPOSE"
	}
	return "ok"
}

/*
**	Create a transaction, pack it to a json string and return the encrypted string.
**  On the client a QR code will be created based on that string.
 */
func GetProposalQRCode(amount int64, purpose string) string {
	trans := _transaction{Amount: amount * -1, Purpose: purpose, Date: time.Now(), Typ: Lmp, From: account.Name, Uuid: account.Uuid}
	jsonData, err := json.Marshal(trans)
	if err != nil {
		log.Println(err)
		return ""
	}
	if purpose == "" {
		log.Println("GetProposalQRCode no purpose")
		return "NO_PURPOSE"
	}
	return encryptStringGCM(string(jsonData), false)
}

/*
**	Create a transation, pack it to a json string and return the encrypted string.
**	On the client a QR code will be created based on that string.
 */
func GetAgreementQRCode() string {
	trans := _transaction{Amount: lastTransaction.Amount * -1, Purpose: lastTransaction.Purpose, Date: time.Now(), Typ: Lmr, From: account.Name, Uuid: lastTransaction.Uuid}
	jsonData, err := json.Marshal(trans)
	if err != nil {
		log.Println(err)
		return ""
	}
	if lastTransaction.Purpose == "" {
		log.Println("GetAgreementQRCode no purpose")
		return "NO_PURPOSE"
	}
	return encryptStringGCM(string(jsonData), false)
}

/*
**	We decrypt the string from a QR code, unpack and validate the transaction
**	and send back the json to the client (we cannot send objects to Kotlin).
**  If it cannot be decrypted or unpacked its a sign for an attack.
 */
func GetProposalFromQRCode(enc string) string {
	jsonData, err := decryptStringGCM(enc)
	if err != nil {
		log.Println("BookTransaction: error decrypting transaction")
		return "FRAUD"
	}
	err = json.Unmarshal([]byte(jsonData), &lastTransaction)
	if err != nil {
		log.Println("BookTransaction: error unmarshaling transaction")
		return "FRAUD"
	}
	if lastTransaction.Typ != Lmp {
		log.Println("BookTransaction: wrong transaction typ")
		return "FRAUD"
	}
	if lastTransaction.Typ == Lmp && lastTransaction.Amount > 0 {
		log.Println("BookTransaction: amount > 0")
		return "FRAUD"
	}
	if lastTransaction.Purpose == "" {
		log.Println("GetProposalFromQRCode no purpose")
		return "NO_PURPOSE"
	}
	return string(jsonData)
}

/*
**	We encrypt the QR-Code, unpack and validate the transaction and when
**	everything is fine we book it.
 */
func GetAgreementFromQRCode(enc string) string {
	jsonData, err := decryptStringGCM(enc)
	if err != nil {
		log.Println("GetAgreementFromQRCode: error decrypting transaction")
		return "FRAUD"
	}
	err = json.Unmarshal([]byte(jsonData), &lastTransaction)
	if err != nil {
		log.Println("GetAgreementFromQRCode: error unmarshaling transaction")
		return "FRAUD"
	}
	if lastTransaction.Typ != Lmr {
		log.Println("GetAgreementFromQRCode: wrong transaction typ")
		return "FRAUD"
	}
	if lastTransaction.Amount < 1 {
		log.Println("GetAgreementFromQRCode: error amount < 1")
		return "WRONG_TYP"
	}
	if lastTransaction.Uuid != account.Uuid {
		log.Println("GetAgreementFromQRCode: error transaction not for this user")
		return "FRAUD"
	}
	if transactionExists(lastTransaction) {
		log.Println("GetAgreementFromQRCode: transaction already booked")
		return "DOUBLE_SPENT"
	}
	if lastTransaction.Purpose == "" {
		log.Println("GetAgreementFromQRCode no purpose")
		return "NO_PURPOSE"
	}
	// only check error on withdraw, not on receive
	addTransaction(lastTransaction.Amount, lastTransaction.Purpose, lastTransaction.Date, lastTransaction.From, lastTransaction.Typ, lastTransaction.Uuid)

	return "ok"
}
