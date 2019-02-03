FROM alpine:3.8

ARG RELEASE_VERSION="0.1.0"

RUN mkdir -p /app && \
    wget -O /app/moni "https://github.com/adrian-gheorghe/moni/releases/download/${RELEASE_VERSION}/moni-linux" && \
    chmod +x /app/moni && \
    cp /app/moni /usr/local/bin/moni