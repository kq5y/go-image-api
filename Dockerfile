FROM golang:1.24-bullseye AS builder
RUN apt-get update && apt-get install -y \
    build-essential \
    pkg-config \
    libvips-dev \
    libheif-dev \
    libheif1 \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o image-api

FROM debian:bullseye-slim
RUN apt-get update && apt-get install -y \
    libvips42 \
    libheif1 \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /root/
COPY --from=builder /app/image-api .
EXPOSE 8080
CMD ["./image-api"]
