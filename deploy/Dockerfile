FROM alpine:20240329

RUN apk add bash

RUN echo "alias ll='ls -al --color'" > ~/.bashrc

RUN mkdir -p /opt/config

RUN mkdir -p /opt/logs

ADD ./clash-ws /opt/clash-ws

ADD ./config /opt/config/

RUN chmod +x /opt/clash-ws

WORKDIR /opt/

ENTRYPOINT ["/opt/clash-ws", "-c", "/opt/config"]
