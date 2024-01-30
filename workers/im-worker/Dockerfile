FROM debian:bookworm-slim

RUN apt update && \
    apt upgrade && \
    apt install imagemagick -y && \
    apt autoremove

COPY ./config/env/ /home/im-worker/config/env
COPY ./im-worker /home/im-worker/im-worker
WORKDIR "/home/im-worker/"

RUN mkdir -p res/tmp \
 res/pending \
 res/processed && \
 ls -la

ENTRYPOINT [ "./im-worker" ]
