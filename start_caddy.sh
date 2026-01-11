#!/usr/bin/env bash
set -euo pipefail

echo Running caddy...

path_to_back="$PWD/apps/caddy"
lsof=$(command -v lsof)

echo "Checking for lsof..."
if ! [[ $lsof ]]; then
    (
        echo "lsof not found"
        read -r -p "Do you want to install lsof? [y/n] " response
        if [[ "$response" =~ ^([yY][eE][sS]|[yY])$ ]]; then
            (
                read -p "Enter install command (with sudo) (i am not doing this for every OS): " cmd
                eval $cmd
            )
        else
            (
                echo "Cannot continue without lsof, exiting..."
                exit 1
            )
        fi
    )
fi

echo "Checking for services on port :80..."
if [[ $(sudo lsof -i:80) ]]; then
    (
        echo "Found service, killing violently..."
        sudo kill -9 $(sudo lsof -t -i:80)
    )
else
    (
        echo "No service found, proceeding..."
    )
fi

if [[ -d "$path_to_back" ]]; then
    (
        cd "$path_to_back"
        sudo caddy run
    )
fi

exit 0