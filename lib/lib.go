package lib

var dbFile string
var secretKey []byte
var account Account
var algorithm string

const (
	algoEnc = "44880922281e9a7d9ba243b2599fcbcf41fd8199b0d059a7b373abdc5f0447380b13b620239871ba3b7db08ca8"
)

func Init(filesDir string) {
	dbFile = filesDir + "/shift.db"
	generateSecretKey()
	algorithm = DecryptStringGCM(algoEnc)
}
