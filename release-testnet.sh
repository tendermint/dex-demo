#!/usr/bin/env bash

commitID=$(git rev-parse HEAD)
branch="testnet"
owner="tendermint"
repo="dex-demo"
info=$(gothub info -u tendermint -r $repo -t $branch)
if [[ $info = *$commitID* ]]; then
    echo "No new commits since last nightly."
    exit 1
fi

gothub delete --user $owner --repo $repo --tag $branch

git tag --force $branch $commitID
git push --force --tags

make clean
make all-cross

mkdir -p ./build/dex-testnet-linux-x64
mv ./build/dexd-linux-x64 ./build/dex-testnet-linux-x64/dexd
mv ./build/dexcli-linux-x64 ./build/dex-testnet-linux-x64/dexcli
tar -zcvf ./build/dex-testnet-linux-x64.tar.gz -C ./build ./dex-testnet-linux-x64

mkdir -p ./build/dex-testnet-darwin-x64
mv ./build/dexd-darwin-x64 ./build/dex-testnet-darwin-x64/dexd
mv ./build/dexcli-darwin-x64 ./build/dex-testnet-darwin-x64/dexcli
tar -zcvf ./build/dex-testnet-darwin-x64.tar.gz -C ./build ./dex-testnet-darwin-x64

echo "Creating release..."
gothub release --user $owner --repo $repo --tag $branch --name "Testnet Build" --description "Manual testnet build." --pre-release

echo "Uploading binaries..."
gothub upload --user $owner --repo $repo --tag $branch --file ./build/dex-testnet-linux-x64.tar.gz --name dex-testnet-linux-x64.tar.gz
gothub upload --user $owner --repo $repo --tag $branch --file ./build/dex-testnet-darwin-x64.tar.gz --name dex-testnet-darwin-x64.tar.gz