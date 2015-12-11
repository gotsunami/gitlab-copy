#!/bin/sh

# Writes a version file for commands that need it.

COMMIT=$(git log --format="%h" -n 1)
TAG=$(git describe --all --exact-match $COMMIT)

for d in src/cmd/*; do
    VERSION="$d/version.go"
    if [ -f $VERSION ]; then
        grep "$COMMIT $TAG" $VERSION >/dev/null || continue
    else
        cat > $VERSION << EOF
package main

const (
    appVersion = "$COMMIT $TAG"
)
EOF
    fi
done
