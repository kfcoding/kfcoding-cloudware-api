FROM ubuntu

MAINTAINER "wsl <wsl@kfcoding.com>"

ADD ./controller /usr/bin/

EXPOSE 8080
ENTRYPOINT ["controller"]