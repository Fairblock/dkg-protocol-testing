import subprocess

def extract_addresses(output):
   
    addresses = []
    lines = output.split('\n')
    skip_next = False
    for line in lines:
        if skip_next:
            skip_next = False
            addresses.append(line.strip())
        if line.startswith('- address:'):
            skip_next = True
    return addresses

def run_export_command(name):

    command = f"~/go/bin/dkgd keys export {name} --unsafe --unarmored-hex"
   # print(command)
    process = subprocess.Popen(command, shell=True, stdin=subprocess.PIPE, stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True)
    output, error = process.communicate(input='y\n')
    #print(output)
    return output

command = "~/go/bin/dkgd keys list"
process = subprocess.Popen(command, shell=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
output, error = process.communicate()

if process.returncode == 0:
    addresses = extract_addresses(output.decode('utf-8'))
   
    private_keys = []
    for address in addresses:
        name = address.split('name:')[1].strip()
        if name != "alice":
            private_key = run_export_command(name)
            private_keys.append(private_key)
    private_keys_string = ','.join(private_keys)
    print(private_keys_string)
else:
    print("An error occurred:", error.decode('utf-8'))
