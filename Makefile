REPO              = github.com/pkar/runit
COMPONENT         = runit
CMD               = $(REPO)/cmd/$(COMPONENT)
IMAGE_NAME        = $(COMPONENT)
IMAGE_TAG         = latest
IMAGE_SPEC        = $(IMAGE_NAME):$(IMAGE_TAG)
UNAME             := $(shell uname)

vendor:
	git remote add -f log git@github.com:pkar/log.git
	git subtree add --squash --prefix=vendor/log log master
	git remote add -f fsnotify git@github.com:go-fsnotify/fsnotify.git
	git subtree add --squash --prefix=vendor/fsnotify fsnotify master
	$(MAKE) vendor_sync

vendor_sync:
	git fetch log
	git subtree pull --message "merge log" --squash --prefix=vendor/log log master
	git fetch fsnotify
	git subtree pull --message "merge fsnotify" --squash --prefix=vendor/fsnotify fsnotify master

build_docker:
	docker build --pull -t $(IMAGE_SPEC) .
	docker run -v $(CURDIR)/bin/linux_amd64:/go/bin $(IMAGE_SPEC) go install $(CMD)

build_linux:
	go get ./...
	mkdir -p bin/linux_amd64
	GOARCH=amd64 GOOS=linux go build -o bin/linux_amd64/$(COMPONENT) ./cmd/$(COMPONENT)/main.go

build_mac:
	go get ./...
	mkdir -p bin/darwin_amd64
	go build -o bin/darwin_amd64/$(COMPONENT) ./cmd/$(COMPONENT)/main.go

build:
ifeq ($(UNAME),Darwin)
	$(MAKE) build_mac
endif
ifeq ($(UNAME),Linux)
	$(MAKE) build_linux
endif

install:
	go install $(CMD)

run:
	go run cmd/$(COMPONENT)/main.go

test:
	go test -cover ./...

testf:
	# make testf TEST=TestRunCmd
	go test -v -test.run="$(TEST)"

bench:
	go test ./... -bench=.

vet:
	go vet ./...

coverprofile:
	# run tests and create coverage profile
	go test -coverprofile=coverage.out .
	# check heatmap
	go tool cover -html=coverage.out

.PHONY: vendor
