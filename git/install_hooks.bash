#!/bin/bash

path="$(dirname $0)"
path="$(cd $path && pwd)"

for file in $(ls "$path/hooks"); do
    chmod +x "$path/hooks/$file" 
    symlink="${file%.*}"
    if [ ! -L ".git/hooks/$symlink" ]; then
        echo "Installing $symlink hook..."
        ln -s "../../git/hooks/$file" ".git/hooks/$symlink"
    fi
done
