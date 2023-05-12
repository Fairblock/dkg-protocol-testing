[All instrcutions and hardcoded values are just for testing purposes and will be replaced with proper formatting.]

To execute the DKG, four repositories are needed and they all need to be located in the same directory. Also please change the repository names according to the below list: 
- DKGChain -> dkg
- DKGCore -> axelar-core
- dkg-protocol -> tofn
- tofnd -> tofnd

This naming convention is temporary and just for testing purposes. 
To run the dkg, perform the following commands from the main directory:
```sh
cd dkg
ignite chain serve
cd ..
cd tofnd
cargo run
cd ..
cd axelar-core
```
Then, based on the two default accound addresses in dkg chain, alice and bob, replace the addresses in line 33-34 of `cmd/axelard/cmd/vald/tss/keygen.go`. Next, run the following commands in two separate terminals to perform the dkg. 
```sh
go run cmd/axelard/main.go vald-start --validator-addr [bob's address] --validator-key [bob's key] 
go run cmd/axelard/main.go vald-start --validator-addr [alice's address] --validator-key [alice's key]
```
