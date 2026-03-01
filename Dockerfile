FROM golang:1.24.3-alpine

WORKDIR /app

RUN apk add --no-cache git bash

COPY go.mod go.sum ./
RUN go mod download

COPY . .

EXPOSE 8082

CMD ["go", "run", "./cmd/server"]
