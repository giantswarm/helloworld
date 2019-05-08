FROM golang:1.11-alpine3.9 AS build

WORKDIR /build

# Compile our helloworld executable
ADD main.go .
RUN go build -o helloworld

# ---------

FROM alpine:3.9

WORKDIR /

# Copy the helloworld executable from the build phase
COPY --from=build /build/helloworld /

# Add our static content
ADD content /content

EXPOSE 8080

ENTRYPOINT ["/helloworld"]
