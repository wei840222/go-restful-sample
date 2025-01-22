FROM golang:1.23.4-bookworm AS builder

WORKDIR /src

COPY go.* ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=1 go build -v -o go-restful-sample

FROM debian:bookworm-slim

# update ca-certificates
RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

ARG user=go-restful-sample
ARG group=go-restful-sample
ARG uid=10000
ARG gid=10001

# If you bind mount a volume from the host or a data container,
# ensure you use the same uid
RUN groupadd -g ${gid} ${group} \
    && useradd -l -u ${uid} -g ${gid} -m -s /bin/bash ${user}

USER ${user}

COPY --from=builder --chown=${uid}:${gid} /src/go-restful-sample /usr/bin/go-restful-sample
COPY --from=builder --chown=${uid}:${gid} /src/config/config.yaml /etc/go-restful-sample/config.yaml

ENV LOG_COLOR=false
ENV LOG_LEVEL=info
ENV LOG_FORMAT=json
ENV GIN_MODE=release

EXPOSE 8080

ENTRYPOINT ["go-restful-sample"]
