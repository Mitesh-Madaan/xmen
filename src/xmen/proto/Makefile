PWD:=$(shell pwd)

all:
	cd ${PWD}
	git clone git@github.com:golang/protobuf.git
	cd protobuf/protoc-gen-go/
	git checkout v1.3.0
	go build
	export PATH=${PWD}/protobuf/protoc-gen-go/:${PATH}
	cd -
	protoc --go_out=. *.proto
