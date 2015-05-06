PROJECT=helloworld
ORGANIZATION=giantswarm
REGISTRY = registry.giantswarm.io
USERNAME := $(shell swarm user)
SOURCE := $(shell find . -name '*.go')
GOPATH := $(shell pwd)/.gobuild
PROJECT_PATH := $(GOPATH)/src/github.com/$(ORGANIZATION)
GOOS := linux
GOARCH := amd64

.PHONY=all clean deps $(PROJECT) swarm-up docker-build docker-push docker-pull

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
	    golang:1.3.1-cross \
	    go build -a -o $(PROJECT)

docker-build: $(PROJECT)
	docker build -t $(REGISTRY)/$(USERNAME)/$(PROJECT) .

docker-push: docker-build
	docker push $(REGISTRY)/$(USERNAME)/$(PROJECT)

docker-pull:
	docker pull $(REGISTRY)/$(USERNAME)/$(PROJECT)

docker-run: docker-build
	docker run -p 8080:8080 -ti --rm $(REGISTRY)/$(USERNAME)/$(PROJECT)

swarm-delete:
	swarm delete $(PROJECT)

swarm-up: docker-push
	swarm up swarm.json --var=username=$(USERNAME)