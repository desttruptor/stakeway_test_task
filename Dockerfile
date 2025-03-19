FROM golang:1.20-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o /validator-api .

FROM alpine:3.17

WORKDIR /app
COPY --from=builder /validator-api .

RUN mkdir -p /data
ENV DB_PATH=/data/validator.db

EXPOSE 8080

CMD ["./validator-api"]