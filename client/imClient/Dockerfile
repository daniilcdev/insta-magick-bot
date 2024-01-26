FROM alpine:3.19

RUN apk update && \
    apk upgrade && \
    apk add --no-cache imagemagick

CMD mkdir -p /home/im-worker/res/tmp \
 /home/im-worker/res/pending \
 /home/im-worker/res/processed

ENTRYPOINT /sbin/init