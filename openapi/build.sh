#!/usr/bin/env bash
set -e

git checkout master
mkdir ./codegen-out
mkdir ./build
./node_modules/.bin/sc generate -l openapi -i ./openapi.yml -o ./codegen-out
git checkout release
git merge master --no-edit
cp ./codegen-out/openapi.json ./build
cp ./dex-demo-embedded.html ./build
git add -f ./build
now=$(date)
git commit -am "Website release $now"
git push origin release
rm -rf ./build
rm -rf ./codegen-out
git checkout master