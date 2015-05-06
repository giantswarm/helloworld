FROM busybox:ubuntu-14.04

ADD ./helloworld /usr/bin/

EXPOSE 8080

ENTRYPOINT ["helloworld"]