FROM golang:1.21 AS go

RUN mkdir app
WORKDIR /app
COPY . .

RUN go build main.go
RUN ls -la helpers
RUN pwd

FROM debian:latest

RUN apt-get update && apt-get install -y ca-certificates && apt-get clean

RUN mkdir /config

WORKDIR /app

COPY --from=go /app/main ./
COPY --from=go /app/helpers/config.toml /config

RUN mkdir /library
RUN mkdir /app/cache
RUN mkdir /app/temp
RUN ls -la /config
RUN pwd
RUN chmod +x main
RUN ls -la /app/main

VOLUME [ "/config", "/library" ]
ENTRYPOINT ["/app/main"]