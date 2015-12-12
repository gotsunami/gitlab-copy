
.PHONY: build dist linux darwin windows buildall version cleardist clean

BIN=gitlab-copy
DISTDIR=dist
GCDIR=${DISTDIR}/${BIN}
#
VERSION=`git describe --tags --always`
GC_VERSION=${BIN}-${VERSION}
GC_DARWIN_AMD64=${GC_VERSION}-darwin-amd64
GC_FREEBSD_AMD64=${GC_VERSION}-freebsd-amd64
GC_LINUX_AMD64=${GC_VERSION}-linux-amd64
GC_WINDOWS_AMD64=${GC_VERSION}-windows-amd64
#
GB_BUILD64=GOARCH=amd64 gb build

all: build

build: version
	@gb build

test: version
	@gb test

coverage:
	@./tools/coverage.sh `pwd`

dist: cleardist buildall zip

zip: linux darwin freebsd windows
	@rm -rf ${GCDIR}

linux:
	@cp bin/${BIN} ${GCDIR}/${BIN} && \
		(cd ${DISTDIR} && zip -r ${GC_LINUX_AMD64}.zip ${BIN})

darwin:
	@cp bin/${BIN}-darwin* ${GCDIR}/${BIN} && \
		(cd ${DISTDIR} && zip -r ${GC_DARWIN_AMD64}.zip ${BIN})

windows:
	@cp bin/${BIN}-windows* ${GCDIR}/${BIN} && \
		(cd ${DISTDIR} && zip -r ${GC_WINDOWS_AMD64}.zip ${BIN})

freebsd:
	@cp bin/${BIN}-freebsd* ${GCDIR}/${BIN} && \
		(cd ${DISTDIR} && zip -r ${GC_FREEBSD_AMD64}.zip ${BIN})

buildall: version
	@GOOS=darwin ${GB_BUILD64}
	@GOOS=freebsd ${GB_BUILD64}
	@GOOS=linux ${GB_BUILD64}
	@GOOS=windows ${GB_BUILD64}

version:
	@./tools/version.sh

cleardist:
	@rm -rf ${DISTDIR} && mkdir -p ${GCDIR}

clean:
	@rm -rf bin pkg ${DISTDIR}
