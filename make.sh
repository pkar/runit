#!/usr/bin/env bash
set -eu

PWD=`pwd`
UNAME=${UNAME:-`uname | tr '[:upper:]' '[:lower:]'`}
ARCH=${ARCH:-"amd64"}
GOOS=${GOOS:-"linux"}
IMAGE_TAG=$(git branch | cut -d ' ' -f 2 | tr -d '\040\011\012\015' | tr "/" "_")
REPO=github.com/pkar/runit
COMPONENT=runit
IMAGE_NAME=pkar/$COMPONENT
IMAGE_SPEC=$IMAGE_NAME:$IMAGE_TAG
TAG=v0.0.6
REPO=github.com/pkar/$COMPONENT
CMD=$REPO/cmd/$COMPONENT
PACKAGE=${PACKAGE:-""}

set -x

path() {
	set +x
	echo export GOPATH=$GOPATH
	echo export PATH=$PATH
}

build() {
	mkdir -p bin/${UNAME}_${ARCH}
	GOARCH=$ARCH GOOS=${UNAME} go build -o bin/${UNAME}_${ARCH}/$COMPONENT-$TAG -ldflags "-X $REPO.AppVersion=${TAG}" `pwd`/cmd/$COMPONENT/main.go
}

install() {
	build
	chmod +x bin/${UNAME}_${ARCH}/$COMPONENT-$TAG
	cp bin/${UNAME}_${ARCH}/$COMPONENT-$TAG /usr/local/bin/$COMPONENT
}

_release() {
	build
	cd bin/${UNAME}_amd64
	tar -czvf runit-${TAG}.${UNAME}.tar.gz runit
	cd -
	mv bin/${UNAME}_amd64/runit-${TAG}.${UNAME}.tar.gz bin/
}

release() {
	UNAME=linux
	_release
	UNAME=darwin
	_release
	UNAME=darwin
}

test() {
	go test -cover .
	#golint .
	go vet .
}

testv() {
	go test -v -cover ./...
}

cover() {
	# run tests and create coverage profile
	go test -coverprofile=coverage.out .
	# check heatmap
	go tool cover -html=coverage.out
}

run() {
	go run src/$REPO/cmd/$COMPONENT/main.go
}

bench() {
	go test ./... -bench=.
}

vendor() {
	glide update --update-vendored --strip-vcs
}

docker_test() {
	docker build -t $IMAGE_SPEC .
	docker run --rm -t $IMAGE_SPEC go test -cover ./...
}

eval $1
