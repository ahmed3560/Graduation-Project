# Stage 1: build the Go binaries
FROM golang:1.27.4-alpine AS build
WORKDIR /src

# Cache dependencies before copying the rest of the source
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/api ./cmd/api && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/worker ./cmd/worker

# Stage 2: small runtime image
FROM alpine:3.20
RUN apk --no-cache add ca-certificates tzdata && \
    adduser -D -u 10001 app

USER app
WORKDIR /app
COPY --from=build /out/ /app/

EXPOSE 8080
CMD ["./api"]
