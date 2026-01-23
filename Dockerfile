# ---- build stage ----
FROM golang:1.22-alpine AS builder
WORKDIR /app

RUN apk add --no-cache git ca-certificates tzdata

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" -o server .

# ---- runtime stage ----
FROM alpine:3.19
WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata \
    && adduser -D -H appuser

COPY --from=builder /app/server /app/server

ENV PORT=8080
EXPOSE 8080
USER appuser

CMD ["/app/server"]
