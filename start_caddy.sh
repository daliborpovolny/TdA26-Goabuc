#!/usr/bin/env bash
set -euo pipefail

echo Running caddy...

path_to_back="$PWD/apps/caddy"

if [[ -d "$path_to_back" ]]; then
    (
        cd "$path_to_back"
        caddy run
    )
fi

exit 0