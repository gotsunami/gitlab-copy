#!/bin/sh

# Writes a version file for commands that need it.

DESC=$(git describe --tags --always)

writeVersion() {
    cat > $1 << EOF
package main

const (
    appVersion = "$DESC"
)
EOF
}

for d in src/cmd/*; do
    VERSION="$d/version.go"
    if [ -f $VERSION ]; then
        grep "\<$DESC\>" $VERSION >/dev/null || writeVersion $VERSION
    else
        writeVersion $VERSION
    fi
done
