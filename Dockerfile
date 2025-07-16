FROM gsoci.azurecr.io/giantswarm/golang:1.24.5-alpine3.21 AS builder

WORKDIR /project

COPY main.go /project/
COPY go.mod /project/
COPY go.sum /project/

RUN go build .

FROM gsoci.azurecr.io/giantswarm/alpine:3.22.1

# Add our static content
ADD content /content

# Add our binary
COPY --from=builder /project/helloworld /helloworld

EXPOSE 8080

ENTRYPOINT ["/helloworld"]
