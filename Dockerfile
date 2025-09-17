FROM golang:1.25.0-alpine

WORKDIR /app

COPY . .

RUN go mod vendor

RUN go build -mod=vendor -o go-clean-hexapp ./cmd/server

CMD ["./go-clean-hexapp"]
