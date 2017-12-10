MAIN_VERSION:=$(shell git describe --abbrev=0 --tags || echo "0.1")
VERSION:=${MAIN_VERSION}\#$(shell git log -n 1 --pretty=format:"%h")

LDFLAGS:=-ldflags "-X inframe/app.Version=${VERSION}"

define \n


endef


default: run


build:
	go build ${LDFLAGS} -a -o server ./cmd/v2

clean:
	del server

test:
	for /f "" %G in ('go list ./... ^| find /i /v "/vendor/"') do @go test %G