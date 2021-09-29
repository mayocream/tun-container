FROM alpine

COPY bin/tunhijack /usr/bin/tunhijack

RUN chmod +x /usr/bin/tunhijack

CMD ["/usr/bin/tunhijack"]