FROM golang:1.25.1-alpine AS build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o server ./cmd/server



FROM alpine:latest
WORKDIR /root/

COPY --from=build /app/server .
COPY .env .

CMD ["./server"]
