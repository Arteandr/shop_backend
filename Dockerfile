FROM golang:1.18-alpine as builder

LABEL stage=gobuilder

ENV GOOS linux

ENV GIN_MODE release

RUN apk update --no-cache && apk add --no-cache tzdata

WORKDIR /build

ADD go.mod .

ADD go.sum .

RUN go mod download

COPY . .

RUN go install github.com/swaggo/swag/cmd/swag@latest

RUN swag init -g ./cmd/app.go

RUN go build -o main ./cmd/app.go

FROM alpine

RUN apk update --no-cache && apk add --no-cache ca-certificates

COPY --from=builder /usr/share/zoneinfo/Europe/Helsinki /usr/share/zoneinfo/Europe/Helsinki

ENV TZ Europe/Helsinki

WORKDIR /app

COPY --from=builder /build/main /app/main
COPY --from=builder /build/configs /app/configs
COPY --from=builder /build/schema /app/schema

EXPOSE 8000
ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.9.0/wait /app/wait
RUN chmod +x /app/wait

CMD ./wait && ./main


