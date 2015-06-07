FROM busybox:ubuntu-14.04

ADD ./helloworld /usr/bin/
ADD content /content

EXPOSE 8080

ENTRYPOINT ["helloworld"]