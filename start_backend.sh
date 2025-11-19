#!/usr/bin/env bash
set -euo pipefail

echo Running backend...

def=$PWD
path_to_main="$def/apps/server/cmd/tourbackend"

if [[ -d "$path_to_main" ]]; then
    (
        cd "$path_to_main"
        go run .
    )
fi

exit 0