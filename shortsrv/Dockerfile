FROM golang:1.16.5-buster as builder
WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o ./out/shortsrv .

FROM alpine:latest
COPY --from=builder /app/out/shortsrv /app/shortsrv
EXPOSE 8081
CMD ["/app/shortsrv"]
