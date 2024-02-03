FROM debian:bookworm-slim

RUN apt-get update && \
    apt-get install -y ca-certificates && \
    apt-get autoremove

WORKDIR /home/telegram-service
RUN mkdir -p ./res/processed

ENTRYPOINT [ "./telegram-service" ]