FROM golang:1.24-alpine AS builder
RUN apk add --no-cache gcc musl-dev vips-dev libheif
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o image-api

FROM alpine:latest
RUN apk add --no-cache vips
WORKDIR /root/
COPY --from=builder /app/image-api .
EXPOSE 8080
CMD ["./image-api"]
