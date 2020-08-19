FROM golang:1.14-alpine AS builder
WORKDIR /go/src/gindriver/
ADD . /go/src/gindriver

# Builder
RUN cd /go/src/gindriver && \
    go build -o gindriver gindriver

# Runner
FROM debian:10

ENV DEBIAN_FRONTEND=noninteractive

# Install mysql & openssh
RUN apt update -y && \
    apt install wget curl ca-certificates -y && \
    apt install mariadb-server -y && \
    apt install openssh-server -y

RUN mkdir /backup && \
    cd /backup && \
    mkdir config && \
    mkdir public

# Fetch binary from builder
COPY --from=builder /go/src/gindriver/gindriver /backup/

# Add static files
COPY public/* /backup/public/

# Entrypoint
COPY entrypoint.sh /entrypoint.sh

RUN chmod 0700 /backup && \
    chmod +x /entrypoint.sh

EXPOSE 3000 22

ENTRYPOINT /entrypoint.sh
