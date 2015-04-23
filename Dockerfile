FROM busybox
ADD index.html index.html
EXPOSE 8080
CMD while true; do nc -l -p 8080 < index.html; done



