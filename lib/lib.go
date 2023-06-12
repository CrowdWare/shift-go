package lib

var dbFile string
var account Account
var algorithm string

const (
	algoEnc = "eb292f2bac576dededa38a1ff0f5280347be2887ffe8ae11d640759dcfe60e6d09140ddb20ae2d3d10a8cb85a6"
)

func Init(filesDir string) {
	dbFile = filesDir + "/shift.db"
	algorithm = DecryptStringGCM(algoEnc)
}
