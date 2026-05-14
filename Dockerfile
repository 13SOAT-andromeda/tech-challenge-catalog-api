FROM golang:1.25-alpine AS production_builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/server ./cmd/api

FROM golang:1.25 AS development
WORKDIR /app
RUN go install github.com/air-verse/air@latest
COPY go.mod go.sum ./
RUN go mod download
COPY . .
CMD ["air", "-c", ".air.toml"]

FROM alpine:3.21 AS production
RUN addgroup -S nonroot && adduser -S nonroot -G nonroot
WORKDIR /app
COPY --from=production_builder /app/server .
USER nonroot:nonroot
EXPOSE 8080
CMD ["./server"]
