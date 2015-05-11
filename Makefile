# Cross compiling a simple helloworld program written in Golang.

PROJECT=helloworld
ORGANIZATION=giantswarm
USERNAME := giantswarm
SOURCE := $(shell find . -name '*.go')
GOPATH := $(shell pwd)/.gobuild
PROJECT_PATH := $(GOPATH)/src/github.com/$(ORGANIZATION)
GOOS := linux
GOARCH := amd64

.PHONY=all clean deps $(PROJECT) docker-build

all: deps $(PROJECT)

clean:
	rm -rf $(GOPATH) $(PROJECT)

deps: .gobuild
.gobuild:
	mkdir -p $(PROJECT_PATH)
	cd $(PROJECT_PATH) && ln -s ../../../.. $(PROJECT)

	# Fetch private packages first (so `go get` skips them later)
	# git clone git@github.com:giantswarm/example.git $(PROJECT_PATH)/example

	# Fetch public packages
	GOPATH=$(GOPATH) go get -d github.com/$(ORGANIZATION)/$(PROJECT)

$(PROJECT): $(SOURCE) 
	echo Building for $(GOOS)/$(GOARCH)
	docker run \
	    --rm \
	    -it \
	    -v $(shell pwd):/usr/code \
	    -e GOPATH=/usr/code/.gobuild \
	    -e GOOS=$(GOOS) \
	    -e GOARCH=$(GOARCH) \
	    -w /usr/code \
	    golang:1.4-cross \
	    go build -a -o $(PROJECT)

docker-build: $(PROJECT)
	docker build -t $(USERNAME)/$(PROJECT) .

docker-run: docker-build
	docker run -p 8080:8080 -ti --rm $(USERNAME)/$(PROJECT)