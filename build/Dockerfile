# Build stage
FROM golang:1.22 as builder

WORKDIR /app
ADD . /app

RUN CGO_ENABLED=0 go build -a -o controller kubernetes/controller/main.go

# Final image
FROM alpine:3.20.2
COPY --from=builder /app/controller /usr/local/bin
WORKDIR /
CMD ["sh"]
