FROM golang:1.22.0

WORKDIR /app

ENV GOOS=linux
ENV CGO_ENABLED=0


# COPY go.mod go.sum ./
RUN mkdir server

COPY ./server/go.mod ./server
COPY ./server/go.sum ./server
RUN cd ./server && go mod download

COPY ./.env ./
COPY ./server ./server

RUN cd ./server && make build

EXPOSE ${PORT}

CMD cd server && ./bin/app
