.PHONY: build dist linux darwin windows buildall cleardist clean

include version.mk

BIN=gitlab-copy
WINDOWS_BIN=${BIN}.exe
DISTDIR=dist
GCDIR=${DISTDIR}/${BIN}
#
GC_VERSION=${BIN}-${VERSION}
GC_DARWIN_AMD64=${GC_VERSION}-darwin-amd64
GC_FREEBSD_AMD64=${GC_VERSION}-freebsd-amd64
GC_OPENBSD_AMD64=${GC_VERSION}-openbsd-amd64
GC_LINUX_AMD64=${GC_VERSION}-linux-amd64
GC_WINDOWS_AMD64=${GC_VERSION}-windows-amd64
#
GB_BUILD64=GOARCH=amd64 go build
MAIN_CMD=github.com/gotsunami/${BIN}/cmd/${BIN}

all: build

build:
	@go build -ldflags "all=$(GO_LDFLAGS)" -o bin/${BIN} ${MAIN_CMD}

test:
	@go test ./... -coverprofile=/tmp/cover.out
	@go tool cover -html=/tmp/cover.out -o /tmp/coverage.html

checksum:
	@for f in ${DISTDIR}/*; do \
		sha256sum $$f > $$f.sha256; \
		sed -i 's,${DISTDIR}/,,' $$f.sha256; \
	done

coverage:
	@./tools/coverage.sh `pwd`

htmlcoverage:
	@./tools/coverage.sh --html `pwd`

dist: cleardist buildall zip sourcearchive checksum

zip: linux darwin freebsd openbsd windows
	@rm -rf ${GCDIR}

linux:
	@cp bin/${GC_VERSION}-linux* ${GCDIR}/${BIN} && \
		(cd ${DISTDIR} && zip -r ${GC_LINUX_AMD64}.zip ${BIN})

darwin:
	@cp bin/${GC_VERSION}-darwin* ${GCDIR}/${BIN} && \
		(cd ${DISTDIR} && zip -r ${GC_DARWIN_AMD64}.zip ${BIN})

windows:
	@cp bin/${GC_VERSION}-windows* ${GCDIR}/${WINDOWS_BIN} && \
		(cd ${DISTDIR} && rm ${BIN}/${BIN} && zip -r ${GC_WINDOWS_AMD64}.zip ${BIN})

freebsd:
	@cp bin/${GC_VERSION}-freebsd* ${GCDIR}/${BIN} && \
		(cd ${DISTDIR} && zip -r ${GC_FREEBSD_AMD64}.zip ${BIN})

openbsd:
	@cp bin/${GC_VERSION}-openbsd* ${GCDIR}/${BIN} && \
		(cd ${DISTDIR} && zip -r ${GC_OPENBSD_AMD64}.zip ${BIN})

buildall:
	@GOOS=darwin ${GB_BUILD64} -v -o bin/${GC_DARWIN_AMD64} ${MAIN_CMD}
	@GOOS=freebsd ${GB_BUILD64} -v -o bin/${GC_FREEBSD_AMD64} ${MAIN_CMD}
	@GOOS=openbsd ${GB_BUILD64} -v -o bin/${GC_OPENBSD_AMD64} ${MAIN_CMD}
	@GOOS=linux ${GB_BUILD64} -v -o bin/${GC_LINUX_AMD64} ${MAIN_CMD}
	@GOOS=windows ${GB_BUILD64} -v -o bin/${GC_WINDOWS_AMD64} ${MAIN_CMD}

sourcearchive:
	@git archive --format=zip -o ${DISTDIR}/${VERSION}.zip ${VERSION}
	@echo ${DISTDIR}/${VERSION}.zip
	@git archive -o ${DISTDIR}/${VERSION}.tar ${VERSION}
	@gzip ${DISTDIR}/${VERSION}.tar
	@echo ${DISTDIR}/${VERSION}.tar.gz

cleardist: clean
	mkdir -p ${GCDIR}

clean:
	@rm -rf bin pkg ${DISTDIR}
