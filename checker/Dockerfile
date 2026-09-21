FROM golang:1.24 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o checker .

FROM alpine:latest
COPY --from=builder /app/checker /checker
CMD ["/checker"]