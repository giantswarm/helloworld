

.PHONY=build

build:
	docker build -t gsoci.azurecr.io/giantswarm/helloworld:latest .

run: build
	docker run -p 8080:8080 -ti --rm gsoci.azurecr.io/giantswarm/helloworld:latest

clean:
	docker rmi gsoci.azurecr.io/giantswarm/helloworld:latest
