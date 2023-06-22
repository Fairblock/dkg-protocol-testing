package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"sync"

	"log"
	"strings"

	"time"

	"math/rand"
	"os"
	"os/exec"

	"github.com/FairBlock/vsskyber"
	bls "github.com/drand/kyber-bls12381"
	"github.com/joho/godotenv"
)

func startChainProgram(id string, threshold string, addressList string, path string) {
	cmd := exec.Command(path, "tx", "dkg", "start-keygen", id, threshold, "1", addressList, "--from", "alice", "-y")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()

	if err != nil {
		fmt.Printf("Failed to start Go program : %s\n", err)
	} else {
		fmt.Printf("Started Go program")
	}
}

func startGoProgram(path string, addr string, key string, port string) {
	cmd := exec.Command("go", "run", "./cmd/dkgd/main.go", "vald-start", "--validator-addr", addr, "--validator-key", key, "--tofnd-port", port)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = path
	err := cmd.Start()

	if err != nil {
		fmt.Printf("Failed to start Go program %s: %s\n", path, err)
	} else {
		fmt.Printf("Started Go program: %s\n", path)
	}
	err = cmd.Wait()
	if err != nil {
		exitErr, ok := err.(*exec.ExitError)
		if ok && exitErr.ExitCode() != 0 {
			fmt.Println("Command exited with non-zero status code:", exitErr.ExitCode())
		} else {
			log.Fatal("Command execution error:", err)
		}
	} else {
		fmt.Println("Command executed successfully")
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Access the environment variables
	addressesString := os.Getenv("ADDRESSES")
	addresses := strings.Split(addressesString, ",")
	keysString := os.Getenv("KEYS")
	keys := strings.Split(keysString, ",")
	portsString := os.Getenv("PORTS")
	ports := strings.Split(portsString, ",")
	shareTestIds := os.Getenv("shareTestIds")

	thresholdString := os.Getenv("THRESHOLD")
	path := os.Getenv("PathTodkgd")
	corePath := os.Getenv("PathToCore")
	if len(addresses) != len(keys) {
		log.Fatal("Mismatch in number of keys and addresses!")
	}
	if len(addresses) != len(ports) {
		log.Fatal("Mismatch in number of ports and addresses!")
	}

	numCalls := len(keys)

	// Create a WaitGroup to wait for goroutines to finish
	var wg sync.WaitGroup

	// Increment the WaitGroup counter for each goroutine
	wg.Add(numCalls)

	for i := 0; i < numCalls; i++ {

		go func(id int) {
			// Decrement the WaitGroup counter when the goroutine finishes
			defer wg.Done()

			// Call the function
			startGoProgram(corePath, addresses[id], keys[id], ports[id])
		}(i)
	}

	jsonBytes, err := json.Marshal(addresses)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	jsonAddressString := string(jsonBytes)
	//~/go/bin/dkgd tx dkg start-keygen 104 1 1 '["cosmos1uvvze65ey932l5l32kfgzlnut8e5f4zp2w26dk","cosmos136wuzlrrceanv5jn0p25um3d426wrc47epsxaj"]' --from alice
	time.Sleep(5 * time.Second)
	rand.Seed(time.Now().UnixNano())

	// Generate a random integer between 0 and 100
	randomNumber := rand.Intn(300)

	startChainProgram(strconv.Itoa(randomNumber), thresholdString, jsonAddressString, path)
	wg.Wait()
	// Keep the programs running until interrupted
	for i := 0; i < len(addresses); i++ {
		copys := corePath + "/share-" + strconv.Itoa(i) + ".txt"
		copyp := corePath + "/pk-" + strconv.Itoa(i) + ".txt"
		//fmt.Println(copy)
		cmd := exec.Command("cp", copys, "./")
		err := cmd.Run()

		if err != nil {
			fmt.Printf("Failed to start Go program: %s\n", err)
		}
		cmd = exec.Command("cp", copyp, "./")
		err = cmd.Run()

		if err != nil {
			fmt.Printf("Failed to start Go program: %s\n", err)
		}
	}
	verify_shares([]string{thresholdString, shareTestIds})

}

func verify_shares(args []string) {
	threshold, _ := strconv.ParseUint(args[0], 10, 32)
	chosen := strings.Split(args[1], ",")

	var shareThreshold []vsskyber.Share

	if len(chosen) < int(threshold) {
		fmt.Println("Number of shares less than threshold.")
		os.Exit(-1)
	}
	for i := 0; i < int(threshold)+1; i++ {

		s, err := ioutil.ReadFile("share-" + string(chosen[i]) + ".txt")
		if err != nil {
			fmt.Printf("Failed to read file: %s\n", err)
			return
		}
		index, _ := strconv.ParseUint(chosen[i], 10, 32)
		fmt.Println(index)
		share := vsskyber.Share{Index: bls.NewKyberScalar().SetInt64(int64(index + 1)), Value: bls.NewKyberScalar().SetBytes(s)}
		shareThreshold = append(shareThreshold, share)
	}

	recMasterSecretKey, err := vsskyber.RegenerateSecret(uint32(threshold)+1, shareThreshold)
	if err != nil{
		fmt.Println(err)
	}
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
		fmt.Println("wrong shares or pk", pk, pkRec, recMasterSecretKey)
		return
	}
}
