FROM golang:1.16.5-buster as builder
WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o ./out/shortsrv ./cmd/shortsrv

FROM alpine:latest
COPY --from=builder /app/out/shortsrv /app/shortsrv
EXPOSE 80
CMD ["/app/shortsrv"]
