import subprocess

def extract_addresses(output):
    addresses = []
    lines = output.split('\n')
    for i in range(len(lines)):
        if lines[i].startswith('- address:'):
            address = lines[i].split('address:')[1].strip()
            name_line = lines[i+1].strip()
           # print(name_line)
            if not name_line.startswith('name: alice'):
                addresses.append(address)
    return addresses

command = "~/go/bin/dkgd keys list"
process = subprocess.Popen(command, shell=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
output, error = process.communicate()

if process.returncode == 0:
    addresses = extract_addresses(output.decode('utf-8'))
    address_string = ','.join(addresses)
    print(address_string)
else:
    print("An error occurred:", error.decode('utf-8'))


