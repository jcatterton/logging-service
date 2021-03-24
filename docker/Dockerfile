FROM golang:1.15.8-alpine as builder
RUN apk update && apk upgrade && apk add --no-cache gcc g++ librdkafka bash libc6-compat git openssh
WORKDIR /login-service
COPY . .
RUN rm -f go.sum
RUN go build -tags musl -o ./app ./cmd/svr/main.go

FROM alpine:3.13.1
WORKDIR /app
COPY --from=builder /login-service/app .
EXPOSE 8004
CMD ["./app"]
