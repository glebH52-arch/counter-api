FROM golang:1.25.0-alpine AS builder

WORKDIR /counter-api

COPY go.sum go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /counter-api/counter-api .

FROM alpine:3.22

RUN addgroup -S counter-api && adduser -S -G counter-api counter-api

WORKDIR /counter-api

COPY --from=builder  /counter-api/counter-api ./counter-api
USER counter-api
EXPOSE 8080
ENTRYPOINT ["./counter-api"]