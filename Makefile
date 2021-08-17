

.PHONY=build

build:
	docker build -t quay.io/giantswarm/helloworld:latest .

run: build
	docker run -p 8080:8080 -ti --rm quay.io/giantswarm/helloworld:latest

clean:
	docker rmi quay.io/giantswarm/helloworld:latest
