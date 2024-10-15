FROM golang:1.22.0

WORKDIR /app

ENV GOOS=linux
ENV CGO_ENABLED=0


# COPY go.mod go.sum ./
RUN mkdir client

COPY ./client/go.mod ./client
COPY ./client/go.sum ./client
RUN cd ./client && go mod download

COPY ./.env ./
COPY ./client ./client

RUN cd ./client && make build

EXPOSE ${PORT}

CMD cd client && ./bin/app
