FROM golang:alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o music-app cmd/app/main.go

# Install golang-migrate CLI
RUN CGO_ENABLED=0 go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.17.0

FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/music-app .
COPY --from=builder /go/bin/migrate ./migrate
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

CMD ["./music-app"]

