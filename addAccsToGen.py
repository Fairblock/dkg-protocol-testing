import subprocess
from dotenv import load_dotenv
import os

# Load environment variables from the .env file
load_dotenv()

# Read the ADDRESSES variable, which should contain a comma-separated list of addresses
addresses = os.getenv('ADDRESSES')

# Split the addresses into a list
address_list = addresses.split(',')

# Iterate through the address list and run the command for each address
for address in address_list:
    command = f"../dkg-chain/dkgd add-genesis-account {address.strip()} 10000000000stake --keyring-backend test"
    subprocess.run(command, shell=True)
