FROM quay.io/giantswarm/golang:1.16.7-alpine3.14 AS builder

WORKDIR /project

COPY main.go /project/
COPY go.mod /project/

RUN go build .

FROM quay.io/giantswarm/alpine:3.14.1

# Add our static content
ADD content /content

# Add our binary
COPY --from=builder /project/helloworld /helloworld

EXPOSE 8080

ENTRYPOINT ["/helloworld"]
