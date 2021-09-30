FROM alpine

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories

RUN apk add --no-cache ca-certificates iptables iproute2 curl

WORKDIR /data

COPY bin/tunhijack /usr/bin/tunhijack

RUN chmod +x /usr/bin/tunhijack

COPY entrypoint.sh entrypoint.sh

RUN chmod +x entrypoint.sh

ENTRYPOINT ["/data/entrypoint.sh"]