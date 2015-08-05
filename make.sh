#!/usr/bin/env bash
set -eu

PWD=`pwd`
UNAME=${UNAME:-`uname | tr '[:upper:]' '[:lower:]'`}
ARCH=${ARCH:-"amd64"}
GOOS=${GOOS:-"linux"}
GOPATH="$PWD/vendor:$PWD"
PATH="$PWD/vendor/bin:$PWD/bin/$ARCH:$PWD:$PATH"
IMAGE_TAG=$(git branch | cut -d ' ' -f 2 | tr -d '\040\011\012\015' | tr "/" "_")
REPO=github.com/pkar/runit
COMPONENT=runit
IMAGE_NAME=pkar/$COMPONENT
IMAGE_SPEC=$IMAGE_NAME:$IMAGE_TAG
TAG=v0.0.2
REPO=github.com/pkar/$COMPONENT
CMD=$REPO/cmd/$COMPONENT
PACKAGE=${PACKAGE:-""}

set -x

path() {
	set +x
	echo export GOPATH="$PWD"
	echo export PATH="$PWD/bin/$ARCH:$PATH"
}

build() {
	mkdir -p bin/${UNAME}_${ARCH}
	GOARCH=$ARCH GOOS=${UNAME} go build -o bin/${UNAME}_${ARCH}/$COMPONENT-$TAG src/$REPO/cmd/$COMPONENT/main.go
}

install() {
	build
	chmod +x bin/${UNAME}_${ARCH}/$COMPONENT-$TAG
	cp bin/${UNAME}_${ARCH}/$COMPONENT-$TAG /usr/local/bin/$COMPONENT
}

release() {
	UNAME=linux
	build
	UNAME=darwin
	build
	cd bin/$UNAME_amd64 && tar -czvf runit-$TAG.$UNAME.tar.gz runit
	mv bin/$UNAME_amd64/runit-$TAG.$UNAME.tar.gz bin/
}

test() {
	go test -cover .
	golint .
	go tool vet --composites=false .
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
	git remote add -f fsnotify git@github.com:go-fsnotify/fsnotify.git
	git subtree add --squash --prefix=vendor/fsnotify fsnotify master
}

vendor_sync() {
	git fetch fsnotify
	git subtree pull --message "merge fsnotify" --squash --prefix=vendor/fsnotify fsnotify master
}

docker_test() {
	docker build -t $IMAGE_SPEC .
	docker run --rm -t $IMAGE_SPEC ./make.sh test
}

eval $1
