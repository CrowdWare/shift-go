package lib

var dbFile string
var account _account
var lastTransaction _transaction

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
