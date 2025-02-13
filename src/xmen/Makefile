PWD:=$(shell pwd)
GOEXEC:=$(shell which go)
PROTOC:=$(shell which protoc)

build_gomike_dev:
	set -ex
	${GOEXEC} version
	cd server/gomike
	${GOEXEC} build -mod=vendor -o gomike
	cp ${PWD}/server/gomike/gomike ${PWD}/../../bin/.
	cd ${PWD}

build_gomike: build_gomike_dev

all: build_gomike

all_dev: build_gomike_dev

refresh_vendor:
	set -ex
	cd server/gomike 
	${GOEXEC} mod tidy
	${GOEXEC} mod vendor
	cd ${PWD}

config:
	$(info ************ MAKE ENV CONFIGURATION ******************)
	$(info GO_EXEC: ${GOEXEC})
	$(info ******************************************************)
	$(info                                                       )

.PHONY:config
