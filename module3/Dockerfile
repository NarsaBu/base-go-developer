FROM golang:1.26-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/server ./cmd/url-shortener

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

RUN adduser -D appuser
USER appuser

WORKDIR /app

COPY --from=builder /app/server .

EXPOSE ${PORT}

CMD ["./server"]