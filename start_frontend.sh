#!/usr/bin/env bash
set -euo pipefail

echo Running frontend...

path_to_back="$PWD/apps/web"

if [[ -d "$path_to_back" ]]; then
    (
        cd "$path_to_back"
        npm run dev -- --host
    )
fi

exit 0
