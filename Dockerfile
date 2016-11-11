FROM alpine:3.4

ADD ./helloworld /usr/bin/
ADD content /content

EXPOSE 8080

ENTRYPOINT ["helloworld"]
