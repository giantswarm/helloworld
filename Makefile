

.PHONY=build

build:
	docker build -t giantswarm/helloworld:latest .

run: build
	docker run -p 8080:8080 -ti --rm giantswarm/helloworld:latest

clean:
	docker rmi giantswarm/helloworld:latest
