FROM golang:1.26-alpine AS builder

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY . .
RUN go build -o /app ./src

FROM alpine:3.21 AS final

RUN mkdir /bin/neuro-schedule-api
COPY --from=builder /app /bin/neuro-schedule-api/app

EXPOSE 80

WORKDIR /bin/neuro-schedule-api

CMD ["./app"]