import subprocess

def extract_addresses(output):
    #print(output)
    addresses = []
    lines = output.split('\n')
    
    j = 0
    for i in range(len(lines)):
        if lines[i].startswith('- address:'):
            address = lines[i].split('address:')[1].strip()
            name_line = lines[i+1].strip()
            print(name_line)
            if not name_line.startswith('name: a1 '):
                if not name_line.startswith('name: v1-key'):
                    addresses.append(address)
                    j = j+1
                # if j==20:
                #     addresses.append("\n\n")
                # if j==88:
                #     break
                #     addresses.append("\n\n")
                
    return addresses

command = "echo 3802380s | ../dkg-chain/dkgd keys list --keyring-backend test"
process = subprocess.Popen(command, shell=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
output, error = process.communicate()

if process.returncode == 0:
    addresses = extract_addresses(output.decode('utf-8'))
    address_string = ','.join(addresses)
    print(address_string)
else:
    print("An error occurred:", error.decode('utf-8'))


