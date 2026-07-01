FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
COPY packages/ ./packages/
COPY services/ ./services/
COPY cmd/ ./cmd/
COPY internal/ ./internal/

RUN go build -o /horizon-api ./cmd/server

FROM alpine:3.19
RUN apk --no-cache add ca-certificates tzdata
COPY --from=builder /horizon-api /usr/local/bin/horizon-api
EXPOSE 8080
CMD ["horizon-api"]
