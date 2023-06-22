package lib

var dbFile string
var peerFile string
var account _account
var lastTransaction _transaction
var initialAmount = int64(initial_amount)
var growLevel0 = int64(10000)
var growLevel1 = int64(1800)
var growLevel2 = int64(360)
var growLevel3 = int64(75)

const (
	algorithm = "AES/GCM/NoPadding"
)

type BalanceError struct {
	message string
}

func (e *BalanceError) Error() string {
	return e.message
}

func Encrypt(value string) string {
	return encryptStringGCM(value, false)
}

func Decrypt(enc string) (string, error) {
	return decryptStringGCM(enc, false)
}

func contains(list []Friend, uuid string) int {
	for index, item := range list {
		if item.Uuid == uuid {
			return index
		}
	}
	return -1
}
