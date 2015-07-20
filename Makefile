REPO              = github.com/pkar/runit
COMPONENT         = runit
CMD               = $(REPO)/cmd/$(COMPONENT)
IMAGE_NAME        = $(COMPONENT)
IMAGE_TAG         = latest
IMAGE_SPEC        = $(IMAGE_NAME):$(IMAGE_TAG)
UNAME             := $(shell uname | awk '{print tolower($0)}')
TAG               = v0.0.2

vendor:
	git remote add -f fsnotify git@github.com:go-fsnotify/fsnotify.git
	git subtree add --squash --prefix=vendor/fsnotify fsnotify master
	$(MAKE) vendor_sync

vendor_sync:
	git fetch fsnotify
	git subtree pull --message "merge fsnotify" --squash --prefix=vendor/fsnotify fsnotify master

build_docker:
	docker build --pull -t $(IMAGE_SPEC) .
	docker run -v $(CURDIR)/bin/linux_amd64:/go/bin $(IMAGE_SPEC) go install $(CMD)

build_linux:
	mkdir -p bin/linux_amd64
	GOARCH=amd64 GOOS=linux go build -o bin/linux_amd64/$(COMPONENT) ./cmd/$(COMPONENT)/main.go

build_darwin:
	mkdir -p bin/darwin_amd64
	go build -o bin/darwin_amd64/$(COMPONENT) ./cmd/$(COMPONENT)/main.go

build:
	$(MAKE) build_$(UNAME)

release:
	$(MAKE) build
	cd bin/$(UNAME)_amd64 && tar -czvf runit-$(TAG).$(UNAME).tar.gz runit
	mv bin/$(UNAME)_amd64/runit-$(TAG).$(UNAME).tar.gz bin/

install:
	go install $(CMD)

run:
	go run cmd/$(COMPONENT)/main.go

test:
	go test -cover .

testv:
	go test -v -cover .

testf:
	# make testf TEST=TestRunCmd
	go test -v -test.run="$(TEST)"

testrace:
	go test -race .

bench:
	go test ./... -bench=.

vet:
	go vet ./...

coverprofile:
	# run tests and create coverage profile
	go test -coverprofile=coverage.out .
	# check heatmap
	go tool cover -html=coverage.out

.PHONY: vendor test install release build
