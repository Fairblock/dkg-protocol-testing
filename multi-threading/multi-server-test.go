package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
)

func main() {
	// Read the .env file
	envFile := ".env"
	namePorts, err := readEnvFile(envFile)
	if err != nil {
		log.Fatal(err)
	}

	// Create a wait group to synchronize goroutines
	var wg sync.WaitGroup

	// Loop through the name-port pairs and run Docker commands in parallel
	for _, np := range namePorts {
		wg.Add(1)
		go func(np string) {
			defer wg.Done()

			// Split the name and port values
			parts := strings.Split(np, ":")
			if len(parts) != 2 {
				log.Printf("Invalid name-port format: %s", np)
				return
			}
			name := parts[0]
			port := parts[1]
			//fmt.Println(name, port)
			// Execute 'sudo docker commit' command
			commitCmd := exec.Command("sudo", "docker", "commit", "tofnd", name)
			
commitCmd.Dir = "/home/ubuntu/tofnd" // Replace with the actual path
commitOutput, commitErr := commitCmd.CombinedOutput()
if commitErr != nil {
   log.Printf("Error committing Docker container for %s: %v\nOutput: %s", name, commitErr, commitOutput)
   return
}
		
runCmd := exec.Command("sudo", "docker", "run", "-p", fmt.Sprintf("%s:%s", port, port), "-i", name, "-p", port)
runCmd.Dir = "/home/ubuntu/tofnd" // Replace with the actual path
runOutput, runErr := runCmd.CombinedOutput()
if runErr != nil {
   log.Printf("Error running Docker container for %s: %v\nOutput: %s", name, runErr, runOutput)
   return
}

outputStr := string(runOutput)
log.Printf("Command output:\n%s", outputStr)
		
		}(np)
	}

	// Wait for all goroutines to finish
	wg.Wait()
}

// readEnvFile reads the .env file and returns a list of name-port pairs
func readEnvFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var namePorts []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		namePorts = append(namePorts, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return namePorts, nil
}
