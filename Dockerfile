FROM busybox

MAINTAINER matthias@giantswarm.io

ADD http-response http-response

EXPOSE 80

# A single command to serve a http response.
# Inspired by http://www.commandlinefu.com/commands/view/9164/one-command-line-web-server-on-port-80-using-nc-netcat
CMD while true; do nc -l -p 80 < http-response; done