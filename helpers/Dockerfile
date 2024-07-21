FROM golang:1.22 AS go

RUN mkdir app
WORKDIR /app
COPY . .

RUN go mod tidy
RUN go install github.com/a-h/templ/cmd/templ@latest
RUN templ generate
RUN go build ./cmd/command/

RUN ls -la helpers
RUN pwd

FROM debian:latest

RUN apt-get update && apt-get install -y ca-certificates && apt-get clean

WORKDIR /app

COPY --from=go /app/command ./

RUN mkdir /data
RUN mkdir /data/temp
RUN mkdir /data/assets

RUN chmod +x command

VOLUME [ "/data", "/library" ]
ENTRYPOINT ["/app/command", "-data-folder", "/data"]