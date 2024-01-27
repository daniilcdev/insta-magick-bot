FROM debian:bookworm-slim

RUN apt update && \
    apt upgrade && \
    apt install imagemagick -y && \
    apt autoremove

RUN mkdir -p /home/im-worker/res/tmp \
 /home/im-worker/res/pending \
 /home/im-worker/res/processed

ENTRYPOINT ["tail" , "-f" , "/dev/null" ]
