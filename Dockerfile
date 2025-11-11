FROM golang:1.25 AS builder

WORKDIR /build
COPY . .

RUN make build

FROM debian:12 AS envTest

RUN apt update && apt install -yq busybox tmux

COPY --from=builder /build/smug /usr/bin
COPY --from=builder /build/completion/smug.bash /etc/bash_completion.d/smug.bash

ENTRYPOINT bash

