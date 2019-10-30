#!/usr/bin/env bash
set -e

rm -rf ./codegen-out
rm -rf ./build
mkdir ./codegen-out
mkdir ./build
./node_modules/.bin/sc generate -l openapi -i ./openapi.yml -o ./codegen-out
cp ./codegen-out/openapi.json ./build
cp ./dex-demo-embedded.html ./build/index.html
git add -f ./build
now=$(date)
git commit -am "Website release $now"
