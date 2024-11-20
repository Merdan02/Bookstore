FROM ubuntu:latest
LABEL authors="merdan"

WORKDIR /Bookstore

COPY go.mod go.sum ./
ENTRYPOINT ["top", "-b"]