FROM golang:1.15-buster AS build
COPY . /src
WORKDIR /src
RUN go build -o app cmd/main.go

FROM debian:buster
RUN apt-get update && \
    apt-get install -y ca-certificates && \
    rm -rf /var/lib/apt/lists/*
COPY --from=build /src/app /app
CMD /app