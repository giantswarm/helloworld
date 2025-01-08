FROM quay.io/giantswarm/golang:1.21.6-alpine3.18 AS builder

WORKDIR /project

COPY main.go /project/
COPY go.mod /project/
COPY go.sum /project/

RUN go build .

FROM quay.io/giantswarm/alpine:3.21.1

# Add our static content
ADD content /content

# Add our binary
COPY --from=builder /project/helloworld /helloworld

EXPOSE 8080

ENTRYPOINT ["/helloworld"]
