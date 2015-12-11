
.PHONY: build dist linux darwin windows buildall cleardist clean

BIN=gitlab-copy
TMPDIR=/tmp/dist
GCDIR=${TMPDIR}/${BIN}
#
GC_DARWIN_AMD64=${BIN}-darwin-amd64
GC_FREEBSD_AMD64=${BIN}-freebsd-amd64
GC_LINUX_AMD64=${BIN}-linux-amd64
GC_WINDOWS_AMD64=${BIN}-windows-amd64
#
GB_BUILD64=GOARCH=amd64 gb build

all: build

build:
	@gb build

test:
	@gb test

coverage:
	@./tools/coverage.sh `pwd`

dist: cleardist buildall zip

zip: linux darwin freebsd windows

linux:
	@cp bin/${BIN} ${GCDIR}/${BIN} && \
		(cd ${TMPDIR} && zip -r ${TMPDIR}/${GC_LINUX_AMD64}.zip ${BIN})

darwin:
	@cp bin/${GC_DARWIN_AMD64} ${GCDIR}/${BIN} && \
		(cd ${TMPDIR} && zip -r ${TMPDIR}/${GC_DARWIN_AMD64}.zip ${BIN})

windows:
	@cp bin/${GC_WINDOWS_AMD64}.exe ${GCDIR}/${BIN} && \
		(cd ${TMPDIR} && zip -r ${TMPDIR}/${GC_WINDOWS_AMD64}.zip ${BIN})

freebsd:
	@cp bin/${GC_FREEBSD_AMD64} ${GCDIR}/${BIN} && \
		(cd ${TMPDIR} && zip -r ${TMPDIR}/${GC_FREEBSD_AMD64}.zip ${BIN})

buildall:
	@GOOS=darwin ${GB_BUILD64}
	@GOOS=freebsd ${GB_BUILD64}
	@GOOS=linux ${GB_BUILD64}
	@GOOS=windows ${GB_BUILD64}

cleardist:
	@rm -rf /tmp/dist && mkdir -p ${GCDIR}

clean:
	@rm -rf bin pkg
	@rm -rf ${TMPDIR}
