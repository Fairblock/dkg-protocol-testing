import subprocess

for i in range(1, 92):
    name = f"a{i}"

    commands = [
        f"../dkg-chain/dkgd keys add {name} --keyring-backend test",
        f"../dkg-chain/dkgd add-genesis-account {name} 100000000stake --keyring-backend test"
    ]
    
    for command in commands:
        try:
            output = subprocess.check_output(command, shell=True, stderr=subprocess.STDOUT)
            print(f"Output for {command}:\n{output.decode()}\n")
        except subprocess.CalledProcessError as e:
            print(f"Error executing command for {command}: {e}\n")

print("Done executing commands.")
