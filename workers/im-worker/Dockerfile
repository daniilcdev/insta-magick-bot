FROM debian:bookworm-slim

RUN apt-get update && \
    apt-get install -y imagemagick ca-certificates && \
    apt-get autoremove

# COPY ./config/env/ /home/im-worker/config/env/
# COPY ./im-worker /home/im-worker/im-worker

# RUN mkdir -p res/tmp \
#  res/pending \
#  res/processed && \
#  ls -la

WORKDIR /home/im-worker/
ENTRYPOINT [ "./im-worker" ]
