FROM golang:1.20-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o app .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
ENV PORT=8080
COPY --from=builder /app/app .
EXPOSE $PORT
ENTRYPOINT ["./app"]
