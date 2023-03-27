FROM alpine:latest
RUN apk add --no-cache netcat-openbsd
COPY check-server.sh /
RUN chmod +x /check-server.sh
ENTRYPOINT ["/bin/sh"]