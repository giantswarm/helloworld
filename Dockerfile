FROM alpine:3.10

WORKDIR /

RUN apk add --no-cache ca-certificates

ADD ./helloworld /helloworld

# Add our static content
ADD content /content

EXPOSE 8080

ENTRYPOINT ["/helloworld"]
