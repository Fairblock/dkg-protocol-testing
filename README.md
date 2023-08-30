# DKG 
The following four repositories provide the complete DKG protocol and its chain:
- dkg-chain 
- dkg-core 
- tofn 
- tofnd 
## Setup 
To run the protocol, the chain and the tofnd instance need to be up and arunning. Each validator runs their own local version of the tofnd. 
To test the DKG, four repositories are needed and they all need to be located in the same directory and tofn should be placed inside the tofnd directory. 

## Testing
First, the chain needs to be started. To start the chain, follow the below commands:
```sh
cd dkg-chain
ignite chain serve
```
Next, the tofnd server should be started. The protocol can be tested either by running a single server for all validators or by simulating the actual case where each validator runs their own local version of the tofnd server. 
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
The ports where the tofnd used by each validator is running on, and final shares that should be used to verify the correctness of the mpk can also be specified in the .env file. For the testing, you also need to specify the timeout for the rounds and the manager key name (who will send the start-keygen tx) in the `.env` file. If the built dkgd file for the chain is not located inside the dkg-chain folder, specify its location in the `.env` file through `PathTodkgd`.
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
### Multiple Servers Testing
In order to run the tests using multiple servers, first we need to setup the chain and run a validator node on each server. Below are the instructions to achieve this:
#### Machine 1:
1. Initialize the chain: `./dkgd init --chain-id <chain_id>`
2. Add new keys with the `keys add` command.
3. Add genesis balance with the `add-genesis-account` command.
4. Run the `gentx` command.

#### Machine 2:
1. Initialize the chain (the chain ID should be the same in both systems):: `./dkgd init --chain-id <chain_id>`
2. Add new keys with the `keys add` command.
3. Add genesis balance with the `add-genesis-account` command.
4. Run the `gentx` command.

#### Machine 1:
1. Copy the `gentx-xxxxxxx*.json` file from machine 2 to the `~/.dkg/confog/gentx` directory. So, the `~/.dkg/confog/gentx` directory in machine 1 should have TWO gentx files; 1 from the original gentx and another from machine 2.
2. Run the same `add-genesis-account` you ran on machine 2. In the command, replace the `key_name` with the `account_address` of the key from machine 2.
3. Run the `collect-gentxs` command.

#### Machine 2:
1. Copy the genesis file from machine 1 and replace the genesis file in machine 2.
2. Add the peer info of machine 1 in the `config.toml` file of machine 2.

#### Both Machines:
- Make sure the 26657 port is accessible. Search for 26657 in the `config.toml` file. If the IP is `127.0.0.1`, change it to `0.0.0.0`.
- Start both the nodes.
If there are more than 2 servers, the same approach can be used to setup the chain and validators.
Next, we can use the previous instructions to run serveral validators on each server for the dkg. Each server will also have their own `tofnd` instance.
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

## DKG Protocol Usage Guide

To effectively utilize the DKG protocol, follow the outlined steps below:

### 1. Set Up Validator Node
First, establish a validator node in the DKG chain. You can do this by following the [node-setup guide](https://github.com/Fairblock/fairyring/blob/main/docs/validator/01_node_setup.md). Make sure to substitute any 'fairyring' specifics with the appropriate DKG chain details.

### 2. Run Local tofnd Server
Once your validator node is up and running, start your own `tofnd` server locally on your machine.

### 3. Initialize DKG Process
- Utilize the `single-server-testing.go` file to begin listening on the DKG chain for the `keygen-start` transaction and continue the protocol through running the `dkg-core`.
    - Make sure to comment the part where the keygen-start messafge is being sent in the `single-server-testing.go` and the part the final shares are verfied, since these parts are only for testing purpose.

### 4. Retrieve Final Share and MPK
- Upon successful completion, the DKG share will be saved in a file named `share-{id}.txt`.
- The Master Public Key (MPK) can be accessed from the chain by querying the `dkg-mpk` event.

This guide should facilitate a smooth setup and operation of the DKG protocol.


## Benchmarking
We have tested the protocol with `5` validators. The number of messages that are being broadcasted in each round are (n is the general number of validators):
- Round 1: 5 broadcast messages (n)
- Round 2: 20 broadcast messages (n * (n-1))
- Round 3: 0 to 20 broadcast messages (depending on the accusations) (0-n*(n-1))
## Number of validators and the required setup
- 20 validators - channelCap = 300 , timeout = 20 blocks
- 30 validators - channelCap = 600 , timeout = 30 blocks
- 50 validators - channelCap = 2500, timeout = 150 blocks
- 180 validators - channelCap = 33000, timeout = 650 blocks

## Overview of the implementation

The dkg implementation is divided into 4 separate projects. Below, we briefly provide a high level description of each part. 

### dkg-chain
The chain side serves as a conduit for broadcasting messages among validators. There are four distinct types of messages. The "start-keygen" message initiates the Distributed Key Generation (DKG) process and includes essential data such as a list of public keys (pks) of the validators participating in the DKG, the round timeout, the session ID, and the threshold value.

Messages within the first two rounds are designated as "refundable messages." During the processing of these messages, an event is emitted, allowing validators to retrieve them. Additionally, these messages are assigned an index during processing. This index serves a critical function; if a validator misses a message, the index enables them to query and retrieve that specific missed message at a later time. A comprehensive list of public keys (pks) is maintained, which will subsequently be utilized in the computation of the Master Public Key (MPK).

The third category of messages is termed a "dispute message." These messages are dispatched solely when a validator intends to allege that another validator is malicious. The validation of the dispute message occurs on the chain. Following the validation process, either the accuser or the accused is declared as faulty through an emitted event. Additionally, a record of all faulty validators is preserved on the chain side. This information is later employed in the calculation of the MPK.

The final message is called keygen-result which is being sent at the end of the protocol and includes the commitments for the pk of each validator. The commitment are later used for the ibe encryption.

### tofn
This section is implemented using the Rust programming language and encompasses the primary components of the Distributed Key Generation (DKG) protocol. The implementation is structured into four distinct rounds:

1. **First Round**: A unique set of keys is generated for each validator, and the corresponding public keys (pks) are broadcasted.
2. **Second Round**: Each validator computes shares for the others and encrypts them using the respective validator's public keys.
3. **Third Round**: Every validator decrypts and verifies the shares they have received. If an incorrect share is identified, a dispute message is created.
4. **Fourth Round**: The final shares are calculated, taking into consideration the faulty validators identified in the previous round.

This structure ensures a methodical and secure progression through the key generation process.

#### Implementation Specific Decisions

Our implementation of the Distributed Key Generation (DKG) protocol is rooted in the methodologies presented in the EthDKG paper. Central to this is the utilization of Zero-Knowledge Proofs (ZKProofs) that the paper introduces for handling anomalous situations, particularly when a validator receives a faulty share and wishes to file a complaint. The proof mechanism is designed to harness a hash function that operates over multiple parameters, such as public keys and the mutual key shared between validators. This allows a validator to convincingly demonstrate their knowledge of their secret key, a vital requirement when alleging the receipt of an erroneous share.

While the original protocol is robust, we've made some modifications to better align it with our technological stack. Specifically, we've adapted the hashing mechanism to suit our use of the bls12-381 elliptic curve. In our implementation, the output of the hash function is converted to a Scalar at a certain point in the algorithm. Therefore, we use the `hash-value mod Modulus` instead of the actual hash value. Importantly, these alterations do not compromise the security integrity of the overall system.

### tofnd
This component acts as a wrapper around the tofn implementation, facilitating communication between the Go side and the Rust side. Its primary role is to manage the continuous stream of messages to and from the dkg-core.

### dkg-core
Serving as the validator code, this section establishes connections to both the chain and the dkg protocol side (tofnd) and mediates the message transfers between them. Several design choices have been made to handle the significant volume of messages, especially when around 200 validators are participating in the DKG, and to ensure synchrony. 

- **Round-Based Messaging**: Messages for each round are only transmitted on-chain during the designated time window for that round, and the corresponding events are read from the chain within the same timeframe. If any messages are missing for a specific round, they will be queried and retrieved from the chain at the beginning of the subsequent round.
- **Batched Transactions**: Messages are sent on-chain in batches, meaning that each transaction dispatched on-chain encompasses a batch of messages.
- **Controlled Delay**: Transactions are sent with a specific delay based on the index of the validator. This strategy ensures that not all messages reach the mempool simultaneously, thereby preventing potential flooding that could lead to missed messages.

## Runtime Measurements
The following section delineates the runtime associated with the worst-case scenario, taking into consideration the varying quantities of participating validators.

| Num of Validators | Num of Malicious Validators | Runtime    |
|-------------------|-----------------------------|------------|
| 10                | 5                           |   4 minutes       |
| 50                | 25                          |   11 minutes      |
| 90                | 45                          |   16 minutes      |
| 140               | 70                          |   0.5 hour    |
| 180               | 90                          |   1 hour      |

