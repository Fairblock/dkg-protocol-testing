# DKG 
The following four repositories provide the complete DKG protocol and its chain:
- dkg-chain 
- dkg-core 
- tofn 
- tofnd 
## Setup 
To run the protocol, the chain and the tofnd instance need to be up and arunning. Each validator runs their own local version of the tofnd. 
To test the DKG, four repositories are needed and they all need to be located in the same directory and tofn should be placed inside the tofnd directory. 

## Usage
First, the chain needs to be started. To start the chain, follow the below commands:
```sh
cd dkg-chain
ignite chain serve
```
Next, the tofnd server should be started. The protocol can be tested either by running a single server for all validators or by simulating the actual case where each validator runs their own local version of the tofnd server using docker. 
### Signle Server Testing
To test the protocol on a signle server, follow the below commands to start the grpc server:
```sh
cd tofnd
cargo build --release
cd target/release
./tofnd
```
The server runs on port 50051 of the localhoast by default. 
The final step is to modify the inputs in .env file and run the single-server-test.go file:
To easily export the addresses for the test, run the ```extract.py``` to extract the addresses, and ```extractkey.py``` to extract the keys.
The ports where the tofnd used by each validator is running on, and final shares that should be used to verify the correctness of the mpk can also be specified in the .env file.
After modifying the inputs, run the test using the below command:
```sh
go run single-server-test.go
```
The test file creates the validators, sends out the key-gen start transaction and performs a validation check on the outputs of the protocol.

### Separate gRPC Servers Inside Docker Containers
To simulate the actual case where each validator has their own version of the gRPC, we can set up the docker containers through the following steps:
```sh
Sudo docker build 
```
Next, create a commit of tofnd for each grpc server required:
```sh
sudo docker commit tofnd tofnd{i}
```
where `i` is the grpc server number `i`. (This is an example of naming.)
Then, run the docker container and map the port on which the gRPC server will be running to a port on the host:
```sh
sudo docker run -p {port_num(e.g. 50051)}:{port_num(e.g. 50051)} -i -t tofnd{i} -p {port_num(e.g. 50051)} 
```
This allows the tofnd to run on ```0.0.0.0:{port_num}``` and be accessible from ```localhost:{port_num}``` on the host.
Create similar containers for all validators but each with a different port on host for the server.
The final step is to modify the inputs for the test (specially the ports) and run the single-server-test.go file:
```sh
go run single-server-test.go
```
The test file creates the validators, sends out the key-gen start transaction and performs a validation check on the outputs of the protocol.
### Malicious Scenarios
To test the protocol in cases where there are faulty validators, use the ```malicious``` feature when building the tofnd:
```sh
cargo build --release --features malicious
```
There are several test cases including:
- A validator sending wrong shares or ciphertext
- A vlidator wrongly accusing someone else for sending wrong shares

In both cases the faulty validator should be excluded from the protocol and the final shares and public key will be created excluding them.
 An example of how to run a malicious case:
```sh
./tofnd malicious R2BadShare 1
```
In this example, anyone running this server will send wrong shares to validator 1.
Other possible options for malicious behaviour are: `R2BadEncryption` and `R3FalseAccusation`
## Benchmarking
We have tested the protocol with `5` validators. The number of messages that are being broadcasted in each round are (n is the general number of validators):
- Round 1: 5 broadcast messages (n)
- Round 2: 20 broadcast messages (n * (n-1))
- Round 3: 0 to 20 broadcast messages (depending on the accusations) (0-n*(n-1))
## Number of validators and the required setup
- 20 validators - channelCap = 300 , timeout = 20 blocks
- 30 validators - channelCap = 600 , timeout = 30 blocks
- 50 validators - channelCap = 2500, timeout = 150 blocks
- 150 Not working yet
