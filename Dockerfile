FROM golang:1.16-alpine AS build
WORKDIR /go/src/git.trap.jp/toki/bot_converter
COPY ./go.* ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 go build -o /converter -ldflags "-s -w"

FROM alpine:3.13.5
WORKDIR /app

RUN apk add --no-cache --update ca-certificates tzdata && \
    update-ca-certificates
ENV DOCKERIZE_VERSION v0.6.1
RUN wget https://github.com/jwilder/dockerize/releases/download/$DOCKERIZE_VERSION/dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && tar -C /usr/local/bin -xzvf dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && rm dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz

EXPOSE 3000

COPY --from=build /converter ./

ENTRYPOINT ["./converter"]
