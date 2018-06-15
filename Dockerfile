FROM ubuntu

MAINTAINER "wsl <wsl@kfcoding.com>"

ADD ./CloudwareController /usr/bin/

EXPOSE 8080
ENTRYPOINT ["CloudwareController"]