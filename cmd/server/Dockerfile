FROM golang:1.13-alpine as builder

RUN mkdir /build
COPY . /build/

WORKDIR /build

RUN go build -o winlose cmd/server/main.go

FROM alpine:latest

COPY --from=builder /build/winlose /app/
WORKDIR /app
RUN chmod +x ./winlose

CMD ["./winlose"]
