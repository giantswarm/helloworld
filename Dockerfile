FROM gsoci.azurecr.io/giantswarm/alpine:3.24.1 AS binary-selector
ARG TARGETPLATFORM
COPY helloworld-* /binaries/
RUN case "$TARGETPLATFORM" in \
      "linux/amd64") cp /binaries/helloworld-linux-amd64 /bin/helloworld ;; \
      "linux/arm64") cp /binaries/helloworld-linux-arm64 /bin/helloworld ;; \
      *) echo "Unsupported platform: $TARGETPLATFORM" && exit 1 ;; \
    esac

FROM gsoci.azurecr.io/giantswarm/alpine:3.24.1

# Add our static content
ADD content /content

# Add our binary
COPY --from=binary-selector /bin/helloworld /helloworld

EXPOSE 8080

ENTRYPOINT ["/helloworld"]
