#!/usr/bin/env bash

if [ "$1" == "debug" ]; then
    set -x
fi

set -euo pipefail

echo -n "Please enter keystore password (input will be hidden): "
# shellcheck disable=SC2034
read -rs KEYSTORE_PASSWORD

echo ""
echo -n "Please enter consensus endpoint: "
# shellcheck disable=SC2034
read -r CONSENSUS_ENDPOINT

echo ""
echo -n "Please enter execution endpoint: "
# shellcheck disable=SC2034
read -r EXECUTION_ENDPOINT

echo ""
echo -n "Please enter wallet withdraw address: "
# shellcheck disable=SC2034
read -r WITHDRAW_ADDRESS

echo ""
echo -n "Please paste in the contents of your keystore file: "
# shellcheck disable=SC2034
read -r KEYSTORE_FILE

mkdir -p /blockchain/validator_keys

echo "$KEYSTORE_FILE" > /blockchain/validator_keys/keys.json

# Set the keystore to be readable by the relay docker user
chown -R 65532:65532 /blockchain/validator_keys/keys.json
chmod 600 /blockchain/validator_keys/keys.json

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
apt-get update

apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin unattended-upgrades apt-listchanges

arch=$(uname -i)
if [[ $arch == arm* ]] || [[ $arch == aarch* ]]; then
    apt-get -y install binfmt-support qemu-user-static
fi


echo unattended-upgrades unattended-upgrades/enable_auto_updates boolean true | debconf-set-selections
dpkg-reconfigure -f noninteractive unattended-upgrades
echo 'Unattended-Upgrade::Automatic-Reboot "false";' >> /etc/apt/apt.conf.d/50unattended-upgrades

docker run --detach \
    --name watchtower \
    --volume /var/run/docker.sock:/var/run/docker.sock \
    containrrr/watchtower

docker run --platform linux/amd64 --detach --name ejector -it -e KEYSTORE_PASSWORD --restart always -v "/blockchain/validator_keys":/keys ghcr.io/vouchrun/pls-lsd-relay:main start \
    --consensus_endpoint "$(echo "$CONSENSUS_ENDPOINT" | xargs)" \
    --execution_endpoint "$(echo "$EXECUTION_ENDPOINT" | xargs)" \
    --keys_dir /keys \
    --withdraw_address "$(echo "$WITHDRAW_ADDRESS" | xargs)"