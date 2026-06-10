FROM golang:alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o music-app cmd/app/main.go

FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/music-app .
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

CMD ["./music-app"]
