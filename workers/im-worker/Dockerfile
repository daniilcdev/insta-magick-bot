FROM debian:bookworm-slim

RUN apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y imagemagick && \
    apt-get autoremove

COPY ./config/env/ /home/im-worker/config/env/
COPY ./im-worker /home/im-worker/im-worker

WORKDIR /home/im-worker/
RUN mkdir -p res/tmp \
 res/pending \
 res/processed && \
 ls -la

ENTRYPOINT [ "./im-worker" ]
