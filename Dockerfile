# Builder
FROM golang:1.17.3-alpine3.14 as builder

RUN apk update && apk upgrade && \
    apk --update add git make

WORKDIR /app

COPY . .

RUN make engine

# Distribution
FROM alpine:latest

RUN apk update && apk upgrade && \
    apk --update --no-cache add tzdata && \
    mkdir /app && mkdir /app/disk

WORKDIR /app 

EXPOSE 4444

COPY --from=builder /app/engine /app

CMD /app/engine