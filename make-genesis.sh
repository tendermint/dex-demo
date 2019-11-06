#!/usr/bin/env bash

rm -rf ~/.dexd*
dexd init testval --chain-id testchain
dexcli keys add dex-demo

dexd add-genesis-account $(dexcli keys show dex-demo -a) 40000000000000000000000000stake,40000000000000000000000000asset1,40000000000000000000000000asset2
dexcli config chain-id testchain
dexcli config output json
dexcli config indent true
dexcli config trust-node true

dexd gentx --name dex-demo

echo "Collecting genesis txs..."
dexd collect-gentxs

echo "Validating genesis file..."
dexd validate-genesis
