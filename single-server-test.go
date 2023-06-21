package main

import (
	"fmt"
	"io"
	"math/rand"
	"strconv"
	"time"

	//	"math/rand"
	"os"
	"os/exec"
	"os/signal"

	// "strconv"
	// "time"

	//"strings"
	"syscall"
)
func startChainProgram(id string) {
	cmd := exec.Command("./bin/dkgd", "tx","dkg","start-keygen",id, "2", "1","[\"cosmos1j58yhcq7atg2re2h6gn6zzgae4s0n979ysxkrz\",\"cosmos1urt3k33qtmnfzumlqvn3d67eulp8ennwvvkd3x\",\"cosmos12hlaf7g85v45433x0d932ctelxvqha6y5nzrsl\",\"cosmos1l2fn2hpdmt2a8sq6697cpm8xha6x8nluanwafx\",\"cosmos1vctx4yuj94k6fk6c7z8tnpeqxhg6h4c9r9hzk2\"]","--from", "alice", "-y")
	
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = "/home/setareh/go"
	err := cmd.Start()
	
	if err != nil {
		fmt.Printf("Failed to start Go program : %s\n", err)
	} else {
		fmt.Printf("Started Go program")
	}
}

func startGoProgram(path string,addr string, key string, port string) {
	cmd := exec.Command("go","run",path, "vald-start", "--validator-addr", addr, "--validator-key", key, "--tofnd-port",port)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = "/home/setareh/job/ibe libs/axelar-core-5aaeee9ca66a9d8cdd11897eb25ad99580ad948b"
	err := cmd.Start()
	
	if err != nil {
		fmt.Printf("Failed to start Go program %s: %s\n", path, err)
	} else {
		fmt.Printf("Started Go program: %s\n", path)
	}
}

func startRustProgram(path string) {

	cmd1 := exec.Command("echo", "3802380")
	cmd2 := exec.Command("./tofnd/target/release/tofnd")
	cmd1.Dir = "/home/setareh/job/ibe libs"
	cmd2.Dir = "/home/setareh/job/ibe libs"
	// Create a pipe to connect the output of cmd1 to the input of cmd2
	pipeReader, pipeWriter := io.Pipe()
	defer pipeReader.Close()
	defer pipeWriter.Close()

	// Set the pipeWriter as the output for cmd1 and pipeReader as the input for cmd2
	cmd1.Stdout = pipeWriter
	cmd2.Stdin = pipeReader
	cmd2.Stdout = os.Stdout
	cmd2.Stderr = os.Stderr
	// Start the commands
	if err := cmd1.Start(); err != nil {
		fmt.Printf("Failed to start command 1: %s\n", err)
		os.Exit(1)
	}
	if err := cmd2.Start(); err != nil {
		fmt.Printf("Failed to start command 2: %s\n", err)
		os.Exit(1)
	}
	//cmd := exec.Command("echo" ,"3802380", "|", "")


}

func main() {

	// Start the Rust program
	//startRustProgram("../tofnd/target/release/tofnd")

	// Start the Go programs 
	startGoProgram("./cmd/dkgd/main.go","cosmos1j58yhcq7atg2re2h6gn6zzgae4s0n979ysxkrz","7bd3c732dfdcf95f8b7308c28eee76a707b655ac89d20caac768867d9764d01a","50059")
	startGoProgram("./cmd/dkgd/main.go", "cosmos1urt3k33qtmnfzumlqvn3d67eulp8ennwvvkd3x", "14d1fffcfaf76b35dfb4fdd1ceb70b75792c276b393b03d366df9c7bd35f4d42","50059")
	startGoProgram("./cmd/dkgd/main.go", "cosmos12hlaf7g85v45433x0d932ctelxvqha6y5nzrsl", "0383dca70f603821066bcd418cfda073cf1dab5a9ccdc8eda4d093d010fbf4da","50059")
	startGoProgram("./cmd/dkgd/main.go", "cosmos1l2fn2hpdmt2a8sq6697cpm8xha6x8nluanwafx", "6db99bd62f9f47b68668c75e9a404605758a799fb4dceb0c0290aa6afbb3402a","50059")
	startGoProgram("./cmd/dkgd/main.go", "cosmos1vctx4yuj94k6fk6c7z8tnpeqxhg6h4c9r9hzk2", "5dc708810b9c285e8a0fd4fafa1e591166ba00793d42c42d7fc5fbdbead9cc91","50059")

	//~/go/bin/dkgd tx dkg start-keygen 104 1 1 '["cosmos1uvvze65ey932l5l32kfgzlnut8e5f4zp2w26dk","cosmos136wuzlrrceanv5jn0p25um3d426wrc47epsxaj"]' --from alice
	time.Sleep(5 * time.Second)
	rand.Seed(time.Now().UnixNano())

	// Generate a random integer between 0 and 100
	randomNumber := rand.Intn(300)

	startChainProgram(strconv.Itoa(randomNumber))
	// Keep the programs running until interrupted
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	fmt.Println("Stopping the programs...")
}
