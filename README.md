# pls-lsd-relay

`pls-lsd-relay` is an off-chain relay service responsible for interacting with PLS Lsd contracts and Beacon chain, synchronizing blocks and events, handling tasks related to validators and balances, and other off-chain operations. Each voter is required to run an instance of this service to submit their votes for new balances of the network, ejected validators and rewards distribution etc.

To learn more about PLS LSD, see [**PLS LSD Documentation and Guide**](https://vouch.run/docs/architecture/vouch_lsd.html)

## Quick Install

On a fresh machine, run the following command to install the relay and setup automatic OS and relay updates:

```bash
curl -sL https://raw.githubusercontent.com/J-tt/pls-lsd-relay/refs/heads/feature/node-install/node-install.sh > node-install.sh; bash node-install.sh
```

Note: this command needs to be run as root.