package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/FairBlock/vsskyber"
	bls "github.com/drand/kyber-bls12381"
)

func reverseBytes(data []byte) {
	length := len(data)
	for i := 0; i < length/2; i++ {
		data[i], data[length-i-1] = data[length-i-1], data[i]
	}
}
func main() {
	args := os.Args
	
	threshold,_ := strconv.ParseUint(args[1], 10, 32)

	var shareThreshold []vsskyber.Share
	if len(args) < int(threshold + 2){
		fmt.Println("Number of shares less than threshold.")
		os.Exit(-1)
	}
	for i := 0; i < int(threshold); i++ {

		s, err := ioutil.ReadFile("share-"+string(args[i+2])+".txt")
	if err != nil {
		fmt.Printf("Failed to read file: %s\n", err)
		return
	}
	index,_ := strconv.ParseUint(args[i+2], 10, 32)
	
	share := vsskyber.Share{Index: bls.NewKyberScalar().SetInt64(int64(index+1)), Value: bls.NewKyberScalar().SetBytes(s)}
	shareThreshold = append(shareThreshold, share)	
}


	recMasterSecretKey, err := vsskyber.RegenerateSecret(uint32(threshold), shareThreshold)

	s := bls.NewBLS12381Suite()
	pkRec := s.G1().Point()
	pkRec.Mul(recMasterSecretKey, s.G1().Point().Base())

	pkb, err := ioutil.ReadFile("pk-0.txt")
	if err != nil {
		fmt.Printf("Failed to read file: %s\n", err)
		return
	}
	pk := s.G1().Point()
	pk.UnmarshalBinary(pkb)

	if pk.Equal(pkRec) {
		fmt.Println("The MSK and reconstructed MSK are equal")
	} else {
		fmt.Println(pk, pkRec, recMasterSecretKey)
		return
	}
}
