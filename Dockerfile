FROM gsoci.azurecr.io/giantswarm/golang:1.26.1-alpine3.23 AS builder

WORKDIR /project

COPY main.go /project/
COPY go.mod /project/
COPY go.sum /project/

RUN go build .

FROM gsoci.azurecr.io/giantswarm/alpine:3.23.3

# Add CA certificates so we can contact external TLS services
RUN apk add --no-cache ca-certificates

# Add our static content
ADD content /content

# Add our binary
COPY --from=builder /project/helloworld /helloworld

EXPOSE 8080

ENTRYPOINT ["/helloworld"]
