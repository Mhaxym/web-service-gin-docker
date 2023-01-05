FROM golang:1.18.1

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -v -o web-service-gin-docker

EXPOSE 8080

CMD ["./web-service-gin-docker"]