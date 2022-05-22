FROM golang:1.17.2

RUN apt-get update
RUN apt-get install -y build-essential git

ENV LANG=C.UTF-8 \
    TZ=Asia/Tokyo \
    GIN_MODE=release \
    TODO_GIN_WORKDIR=/app

WORKDIR /app

COPY . /app
RUN go mod download
RUN go build -o app
