FROM golang:1.14-buster AS builder
WORKDIR /go/src/gindriver/
ADD . /go/src/gindriver

# Builder
RUN cd /go/src/gindriver && \
    CGO_ENABLED=0 go build -o gindriver \
        -a -ldflags '-s -w' gindriver

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
    chmod 0700 /entrypoint.sh

EXPOSE 3000 22

ENTRYPOINT /entrypoint.sh
