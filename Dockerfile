FROM ubuntu
Maintainer "wsl <wsl@kfcoding.com>"

ADD ./KfcodingIngressController /usr/bin/

EXPOSE 8080
ENTRYPOINT ["KfcodingIngressController"]