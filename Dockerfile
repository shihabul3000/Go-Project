FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /spotsync-api .

FROM alpine:3.20

RUN adduser -D -g '' spotsync
USER spotsync

WORKDIR /app
COPY --from=builder /spotsync-api /app/spotsync-api

EXPOSE 8080
CMD ["/app/spotsync-api"]
