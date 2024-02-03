FROM debian:bookworm-slim

RUN apt-get update && \
    apt-get install -y imagemagick ca-certificates && \
    apt-get autoremove

WORKDIR /home/im-worker/
ENTRYPOINT [ "./im-worker" ]
