#!/bin/sh

# Writes a version file for commands that need it.

DESC=$(git describe --always)

for d in src/cmd/*; do
    VERSION="$d/version.go"
    if [ -f $VERSION ]; then
        grep "$DESC" $VERSION >/dev/null || continue
    else
        cat > $VERSION << EOF
package main

const (
    appVersion = "$DESC"
)
EOF
    fi
done
