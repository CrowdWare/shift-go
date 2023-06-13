package lib

var dbFile string
var account Account

const (
	algorithm = "AES/GCM/NoPadding"
)

func Init(filesDir string) {
	dbFile = filesDir + "/shift.db"
}
