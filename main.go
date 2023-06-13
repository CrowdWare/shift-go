package main

import (
	"fmt"

	"github.com/crowdware/shift-go/lib"
)

func main() {
	fmt.Println("Shift")
	enc := lib.EncryptStringGCM("this is a test lorem ipsum dolor simit aroiun ghs", false)
	fmt.Println("enc:", enc)
	plain := lib.DecryptStringGCM(enc)
	fmt.Println("plain:", plain)
	api_key_enc := lib.EncryptStringGCM("1234567890123456"+"9878623447645687", false)
	fmt.Println(api_key_enc)
	servive_url_enc := lib.EncryptStringGCM("http://shift.crowdware.at:8080/", false)
	fmt.Println(servive_url_enc)
}
