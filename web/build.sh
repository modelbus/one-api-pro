#!/bin/sh

while IFS= read -r theme; do
    echo "Building theme: $theme"
    rm -rf build/"$theme"
    cd "$theme"
    npm install
    npm run build
    cd ..
done < THEMES
