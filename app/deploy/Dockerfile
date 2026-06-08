FROM golang:1.26.4-alpine AS builder

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download && go mod verify

COPY . .

ENV CGO_ENABLED=0 GOOS=linux
RUN go build -o main ./cmd/server

FROM gcr.io/distroless/static:nonroot

COPY --from=builder /app/main /usr/local/bin/

USER nonroot

EXPOSE 8080

CMD ["/usr/local/bin/main"]