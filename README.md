# Sentinel CLI Client

[![Go](https://img.shields.io/github/go-mod/go-version/sentinel-official/cli-client)]()
[![GoReport](https://goreportcard.com/badge/github.com/sentinel-official/cli-client)](https://goreportcard.com/report/github.com/sentinel-official/cli-client)
[![Licence](https://img.shields.io/github/license/sentinel-official/cli-client.svg)](https://github.com/sentinel-official/cli-client/blob/master/LICENSE)
[![Tag](https://img.shields.io/github/tag/sentinel-official/cli-client.svg)](https://github.com/sentinel-official/cli-client/releases/latest)
[![TotalLines](https://img.shields.io/tokei/lines/github/sentinel-official/cli-client)]()

## Install dependencies

### Linux

```sh
sudo apt-get update && \
sudo apt-get install curl openresolv wireguard-tools && \
sudo sh -c "curl -fsLS https://raw.githubusercontent.com/v2fly/fhs-install-v2ray/master/install-release.sh | bash -s -- --version v5.2.1"
```

### Mac

```sh
brew install v2ray wireguard-tools
```

or

```sh
port install v2ray wireguard-tools
```

## Install Sentinel CLI client

```sh
curl --silent https://raw.githubusercontent.com/sentinel-official/cli-client/master/scripts/install.sh | sh
```

## Connect to a dVPN node

1. Create or recover a key
   
    Need not perform this step again in case you have already done it once.
   
   ```sh
   sentinelcli keys add \
       --home "${HOME}/.sentinelcli" \
       --keyring-backend file \
       <KEY_NAME>
   ```
   
    Pass flag `--recover` to recover the key.

2. Query the active nodes and choose one
   
   ```sh
   sentinelcli query nodes \
       --home "${HOME}/.sentinelcli" \
       --node https://rpc.sentinel.co:443 \
       --status Active \
       --page 1
   ```
   
    Increase the page number to get more nodes

3. Subscribe to a node
   
   ```sh
   sentinelcli tx subscription subscribe-to-node \
       --home "${HOME}/.sentinelcli" \
       --keyring-backend file \
       --chain-id sentinelhub-2 \
       --node https://rpc.sentinel.co:443 \
       --gas-prices 0.1udvpn \
       --from <KEY_NAME> <NODE_ADDRESS> <DEPOSIT>
   ```

4. Query the active subscriptions of your account address
   
   ```sh
   sentinelcli query subscriptions \
       --home "${HOME}/.sentinelcli" \
       --node https://rpc.sentinel.co:443 \
       --status Active \
       --page 1 \
       --address <ACCOUNT_ADDRESS>
   ```

5. Connect
   
   ```sh
   sudo sentinelcli connect \
       --home "${HOME}/.sentinelcli" \
       --keyring-backend file \
       --chain-id sentinelhub-2 \
       --node https://rpc.sentinel.co:443 \
       --gas-prices 0.1udvpn \
       --yes \
       --from <KEY_NAME> <SUBSCRIPTION_ID> <NODE_ADDRESS>
   ```

## Disconnect from a dVPN node

1. Disconnect
   
   ```sh
   sudo sentinelcli disconnect \
       --home "${HOME}/.sentinelcli"
   ```

Click [here](https://docs.sentinel.co/sentinel-cli "here") to know more!
