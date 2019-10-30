#!/usr/bin/env bash

set -e

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
cd $DIR/ui
npm install
npm run build
find ./build -type f -name "*.map" -exec rm '{}' \;
gsed -i 's+</head>+<script type="text/javascript">window.CSRF_TOKEN=\x27{{.CSRFToken}}\x27;window.UEX_ADDRESS=\x27{{.UEXAddress}}\x27;</script></head>+g' ./build/index.html
rm -rf "$DIR/embedded/ui/public"
mkdir -p "$DIR/embedded/ui/public"
cp -r ./build/* "$DIR/embedded/ui/public"
cd ../