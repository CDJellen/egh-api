FROM golang:1.19 AS build

WORKDIR /build

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN cd cmd/server/ && go build -o egh-api 

EXPOSE 8080

ENTRYPOINT [ "cmd/server/egh-api" ]
