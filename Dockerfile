FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /payflow-api ./cmd/api

FROM gcr.io/distroless/base-debian12
COPY --from=builder /payflow-api /payflow-api
EXPOSE 8080
ENTRYPOINT ["/payflow-api"]
