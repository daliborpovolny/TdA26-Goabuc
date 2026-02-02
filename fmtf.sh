#!/usr/bin/env bash
set -euo pipefail

web="$PWD/apps/web"

if [[ -d "$web" ]]; then
    cd $web
    npm run format
fi

cd ../..

exit 0