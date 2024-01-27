FROM debian:bookworm-slim

RUN apt update && \
    apt upgrade && \
    apt install imagemagick -y && \
    apt autoremove

RUN mkdir -p /home/im-worker/res/tmp \
 /home/im-worker/res/pending \
 /home/im-worker/res/processed

COPY ./env /home/im-worker/env
COPY ./im-worker /home/im-worker/im-worker

WORKDIR "/home/im-worker/"
ENTRYPOINT [ "./im-worker" ]
