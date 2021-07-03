FROM golang:1.16.5-buster as builder
WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o ./out/shortgen .

FROM alpine:latest
COPY --from=builder /app/out/shortgen /app/shortgen
CMD ["/app/shortgen"]
