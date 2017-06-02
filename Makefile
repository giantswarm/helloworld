# Cross compiling a simple helloworld program written in Golang.

PROJECT=helloworld
ORGANIZATION=giantswarm
USERNAME := giantswarm
SOURCE := $(shell find . -name '*.go')
GOOS := linux
GOARCH := amd64

.PHONY=all clean $(PROJECT) docker-build

all: $(PROJECT)

clean:
	rm -rf $(PROJECT)

$(PROJECT): $(SOURCE)
	echo Building for $(GOOS)/$(GOARCH)
	docker run \
		--rm \
		-it \
		-v $(shell pwd):/usr/code \
		-e GOPATH=/usr/code/.gobuild \
		-e GOOS=$(GOOS) \
		-e GOARCH=$(GOARCH) \
		-e CGO_ENABLED=0 \
		-w /usr/code \
		golang:1.8.3 \
		go build -a -installsuffix cgo -o $(PROJECT)

docker-build: $(PROJECT)
	docker build -t $(USERNAME)/$(PROJECT) .

docker-run: docker-build
	docker run -p 8080:8080 -ti --rm $(USERNAME)/$(PROJECT)
