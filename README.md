# pls-lsd-relay

`pls-lsd-relay` is an off-chain relay service responsible for interacting with PLS Lsd contracts and Beacon chain, synchronizing blocks and events, handling tasks related to validators and balances, and other off-chain operations. Each voter is required to run an instance of this service to submit their votes for new balances of the network, ejected validators and rewards distribution etc.

To learn more about PLS LSD, see [**PLS LSD Documentation and Guide**](https://vouch.run/docs/architecture/vouch_lsd.html)

## Quick Install

The below command will:
- Install Docker (if not already installed)
- Configure relay directory and config.toml file
- Import Relay account private keys
- Enable automatic OS and relay updates (optional)
- Start the Relay docker container


Note: this command needs to be run as root.

```bash
curl -sL https://raw.githubusercontent.com/Vouchrun/pls-lsd-relay/refs/heads/main/relay-install.sh > relay-install.sh; sudo bash relay-install.sh
```

## Detailed Install

More information and installation options can be found at the [**Relay Client Documentation**](https://vouch.run/docs/governance/relay_client.html)