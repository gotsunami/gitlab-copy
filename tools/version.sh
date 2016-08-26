#!/bin/sh

# Writes a version file for commands that need it.

VERSION="$(git describe --tags --always)"
GITREV=$(git rev-parse --verify --short HEAD)
GITBRANCH="$(git rev-parse --abbrev-ref HEAD)"
BUILT=$(LANG=US date +"%a, %d %b %Y %X %z")

writeVersion() {
    cat > $1 << EOF
package main

// THIS FILE IS GENERATED AUTOMATICALLY BY tools/version.sh

const (
    Version = "$VERSION"
    GitRevision = "$GITREV"
    GitBranch = "$GITBRANCH"
    Built = "$BUILT"
)
EOF
}

FILE="src/command/gitlab-copy/version.go"
writeVersion $FILE
