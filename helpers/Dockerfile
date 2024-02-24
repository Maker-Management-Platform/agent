FROM golang:1.21 AS go

RUN mkdir app
WORKDIR /app
COPY . .

RUN go build main.go
RUN ls -la helpers
RUN pwd

FROM debian:latest

RUN apt-get update && apt-get install -y ca-certificates && apt-get clean

WORKDIR /app

COPY --from=go /app/main ./

RUN mkdir /data
RUN mkdir /data/temp

RUN chmod +x main

VOLUME [ "/data", "/library" ]
ENTRYPOINT ["/app/main"]