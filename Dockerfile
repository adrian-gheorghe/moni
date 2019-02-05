FROM alpine:3.8

COPY VERSION /app/VERSION

RUN apk add --no-cache curl && \
    RELEASE_VERSION=$(cat /app/VERSION) && \
    echo $RELEASE_VERSION && \
    mkdir -p /app && \
    wget -O /app/moni "https://github.com/adrian-gheorghe/moni/releases/download/${RELEASE_VERSION}/moni-linux" && \
    chmod +x /app/moni && \
    cp /app/moni /usr/local/bin/moni

COPY sample.docker.config.yml /app/config.yml
ENV CONFIG_PATH /app/config.yml

CMD ["sh", "-c", "moni --config ${CONFIG_PATH}"]