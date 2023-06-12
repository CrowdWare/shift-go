package main

import (
	"fmt"

	"github.com/crowdware/shift-go/lib"
)

func main() {
	fmt.Println("Shift")
	enc := lib.EncryptStringGCM("this is a test lorem ipsum dolor simit aroiun ghs")
	fmt.Println("enc:", enc)
	plain := lib.DecryptStringGCM(enc)
	fmt.Println("plain:", plain)
}
