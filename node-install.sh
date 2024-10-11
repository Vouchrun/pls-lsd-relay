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

echo "account = \"$ACCOUNT\"
trustNodeDepositAmount     = 1000000  # PLS
eth2EffectiveBalance       = 32000000 # PLS
maxPartialWithdrawalAmount = 8000000  # PLS
gasLimit = \"3000000\"
maxGasPrice = \"1200\"                            #Gwei
batchRequestBlocksNumber = 16
runForEntrustedLsdNetwork = false

[pinata]
apikey = \"$PINATA\"
pinDays = 180
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
lsdTokenAddress = \"0xLSD_TOKEN_ADDRESS\"
lsdFactoryAddress = \"0xLSD_FACTORY_ADDRESS\"

[[endpoints]]
eth1 = \"https://rpc-pulsechain.g4mm4.io\"
eth2 = \"https://rpc-pulsechain.g4mm4.io/beacon-api/\"
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



read -r -p "Enable automatic updates? [y/n] (press Enter to confirm): " AUTOMATIC_UPDATES

if [[ ${AUTOMATIC_UPDATES:0:1} =~ ^[Yy]$ ]]; then
    echo unattended-upgrades unattended-upgrades/enable_auto_updates boolean true | debconf-set-selections
    dpkg-reconfigure -f noninteractive unattended-upgrades
    echo 'Unattended-Upgrade::Automatic-Reboot "true";' >> /etc/apt/apt.conf.d/50unattended-upgrades

    if ! docker ps -q -f name=watchtower | grep -q .; then
        docker run --detach \
            --name watchtower \
            --volume /var/run/docker.sock:/var/run/docker.sock \
            containrrr/watchtower
    else
        echo "Watchtower is already running, skipping setup..."
    fi
fi

echo ""
# echo -n "Would you like to import the Private Key for your selected Relay Account? [y/n]: "
# read -r -n 1 IMPORT_KEY




read -r -p "Would you like to import the Private Key for your selected Relay Account? [y/n] (press Enter to confirm): " IMPORT_KEY


if [[ $IMPORT_KEY =~ ^[Yy]$ ]]; then
    docker run --name relay-key-import -it -e KEYSTORE_PASSWORD -v "$CONFIG_PATH":/keys ghcr.io/vouchrun/pls-lsd-relay:main import-account --base-path /keys
    docker stop relay-key-import
    docker rm relay-key-import
fi




echo ""
echo -n "Start relay service? (starting now will pass through your key password) [y/n]: "
read -r -n 1 START_SERVICE
echo ""

echo -n "Enter a customised container name for the relay service (default)[relay]: "
read -r RELAY_CONTAINER_NAME

if [ -z "${RELAY_CONTAINER_NAME}" ]; then
    RELAY_CONTAINER_NAME="relay"
fi

if [[ $START_SERVICE =~ ^[Yy]$ ]]; then
    docker run --name "$RELAY_CONTAINER_NAME" -d -e KEYSTORE_PASSWORD --restart always -v "$CONFIG_PATH":/keys ghcr.io/vouchrun/pls-lsd-relay:main start --base-path /keys
else
    echo ""
    echo ""
    echo "To start the relay client, run: "
    echo "docker run --name \"$RELAY_CONTAINER_NAME\" -d  --restart always -v \"$CONFIG_PATH\":/keys ghcr.io/vouchrun/pls-lsd-relay:main start --base-path /keys"
fi