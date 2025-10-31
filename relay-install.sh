#!/usr/bin/env bash

if [ "$1" == "debug" ]; then
    set -x
fi

if [ "$EUID" -ne 0 ]
  then echo "Please run as root"
  exit
fi

set -euo pipefail

echo ""
echo -n "Please enter your voter account address: "
read -r ACCOUNT

echo ""
echo -n "Please enter your Pinata Secret Access Token: "
read -r PINATA

echo ""
echo -n "Please enter your keystore password (your keys will be imported later, your password won't be displayed as you type): "
# shellcheck disable=SC2034
read -r -s KEYSTORE_PASSWORD
echo

echo ""
read -r -p "Use testnet settings [y/n] (press Enter to confirm): " TESTNET

if [[ ${TESTNET:0:1} =~ ^[Yy]$ ]]; then
    DEFAULT_CONFIG_PATH="/blockchain/relay/testnet"
else
    DEFAULT_CONFIG_PATH="/blockchain/relay"
fi

echo ""
echo -n "If you require a customised relay location enter the path (default)[$DEFAULT_CONFIG_PATH]: "
read -r CONFIG_PATH

if [ -z "${CONFIG_PATH}" ]; then
    CONFIG_PATH="$DEFAULT_CONFIG_PATH"
fi

mkdir -p "$CONFIG_PATH"

echo "account = \"$ACCOUNT\"          # Address in checksum format (i.e. case-sensitive)
trustNodeDepositAmount     = 1000000  # PLS
eth2EffectiveBalance       = 32000000 # PLS
maxPartialWithdrawalAmount = 8000000  # PLS
gasLimit = \"3000000\"
maxGasPrice = \"40000000\"              # Gwei
gasPriceMultiplier = 2.5
batchRequestBlocksNumber = 16
eventFilterMaxSpanBlocks = 100000
maxEjectedValPerCycle  = 100          # 0 for unlimited
runForEntrustedLsdNetwork = false
transferFeeAddresses = [\"0x54da21340773FeCAF9A5bad0883a7fc594945d0A\"] # whitelist DAO distributor address
distributeBlockedTransferFeePerEra = 50000000 # for clearing excess fees at VOUCH launch

[pinata]
apikey = \"$PINATA\"
pinDays = 60
" > "$CONFIG_PATH/config.toml"

if [[ ${TESTNET:0:1} =~ ^[Yy]$ ]]
then
    echo "[contracts]
lsdTokenAddress = \"0x61135C59A4Eb452b89963188eD6B6a7487049764\"
lsdFactoryAddress = \"0x98f51f52A8FeE5a469d1910ff1F00A3D333bc9A6\"

[[endpoints]]
eth1 = \"https://rpc-testnet-pulsechain.g4mm4.io\"
eth2 = \"https://rpc-testnet-pulsechain.g4mm4.io/beacon-api/\"
" >> "$CONFIG_PATH/config.toml"
else
    echo "[contracts]
lsdTokenAddress = \"0x79BB3A0Ee435f957ce4f54eE8c3CFADc7278da0C\"
lsdFactoryAddress = \"0x4bf4df49f8bc72a4e484443a14b827cb8c47c716\"

# Primary Local Endpoints
[[endpoints]]
eth1 = \"http://localhost:8545\"
eth2 = \"http://localhost:5052\"

# Vouch RPC Enpoints
[[endpoints]]
eth1 = \"https://rpc.vouch.run\"
eth2 = \"https://rpc-beacon.vouch.run\"

# Fallback G4mm4.io public endpoints
[[endpoints]]
eth1 = \"https://rpc-pulsechain.g4mm4.io\"
eth2 = \"https://rpc-pulsechain.g4mm4.io/beacon-api\"
" >> "$CONFIG_PATH/config.toml"
fi

echo ""
echo "Created default config.toml"

# Set the keystore to be readable by the relay docker user
chown -R 65532:65532 "$CONFIG_PATH"

if docker --version ; then
    echo "Detected existing docker installation, skipping install..."
else
    # Add Docker's official GPG key:
    apt-get update
    apt-get install ca-certificates curl
    install -m 0755 -d /etc/apt/keyrings
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
    chmod a+r /etc/apt/keyrings/docker.asc

    # Add the repository to Apt sources:
    echo \
    "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu \
    $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
    tee /etc/apt/sources.list.d/docker.list > /dev/null
    apt-get -qq update

    apt-get -qq install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin unattended-upgrades apt-listchanges

    arch=$(uname -i)
    if [[ $arch == arm* ]] || [[ $arch == aarch* ]]; then
        apt-get -qq -y install binfmt-support qemu-user-static
    fi
fi


read -r -p "Would you like to import the Private Key for your selected Relay Account? [y/n] (press Enter to confirm): " IMPORT_KEY


if [[ $IMPORT_KEY =~ ^[Yy]$ ]]; then
    docker container prune -f
    docker run --name relay-key-import -it -e KEYSTORE_PASSWORD="$KEYSTORE_PASSWORD" -v "$CONFIG_PATH":/keys ghcr.io/vouchrun/pls-lsd-relay:main import-account --base-path /keys
    docker stop relay-key-import
    docker container prune -f
fi

echo ""
echo " Creating Startup Script"
echo ""
echo -n "Enter a customised container name for the relay service (default)[relay]: "
read -r RELAY_CONTAINER_NAME

if [ -z "${RELAY_CONTAINER_NAME}" ]; then
    RELAY_CONTAINER_NAME="relay"
fi

echo "#!/bin/bash

sudo docker stop $RELAY_CONTAINER_NAME
sudo docker rm $RELAY_CONTAINER_NAME
sudo docker run -it --pull always \
--name $RELAY_CONTAINER_NAME \
-e KEYSTORE_PASSWORD \
--restart unless-stopped \
--network=host \
-v \"$CONFIG_PATH\":/keys \
ghcr.io/vouchrun/pls-lsd-relay:main \
start --base-path /keys \
--log-level info" > "$CONFIG_PATH/start_relay.sh"

chmod +x "$CONFIG_PATH/start_relay.sh"

echo ""
echo " Startup Script Creation Successful"
echo ""


echo ""
read -r -p "Start relay service? [y/n] (press Enter to confirm): " START_SERVICE

if [[ ${START_SERVICE:0:1} =~ ^[Yy]$ ]]; then
    "$CONFIG_PATH/start_relay.sh"
else
    echo ""
    echo ""
    echo "To start the relay client, run: "
    echo "$CONFIG_PATH/start_relay.sh"
fi