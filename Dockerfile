FROM golang:1.18.0-alpine


RUN apk add alpine-sdk

RUN mkdir /app
WORKDIR /app
COPY ./ .

RUN go clean --modcache
RUN apk update && apk upgrade

RUN apk add libressl-dev
RUN apk add openssl-dev

RUN go mod tidy

ENTRYPOINT ["sh"]