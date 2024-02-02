FROM debian:bookworm-slim

RUN apt-get update && \
    apt-get install -y imagemagick ca-certificates && \
    apt-get autoremove

COPY ./config/env/imagemagick.env /home/im-worker/config/env/imagemagick.env
COPY ./im-worker /home/im-worker/im-worker

RUN mkdir -p /home/im-worker/res/tmp \
 /home/im-worker/res/pending \
 /home/im-worker/res/processed && \
 ls -la

WORKDIR /home/im-worker/
ENTRYPOINT [ "./im-worker" ]
