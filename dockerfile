FROM golang:1.25.6-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o binary
FROM scratch
COPY --from=builder /app/binary /binary
EXPOSE 8080
ENTRYPOINT ["/binary"]