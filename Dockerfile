# docker build . -t uzrnem/expense:0.2.1 --no-cache
FROM golang:1.21rc2-alpine3.18 AS build
WORKDIR /temp
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main ./cmd/app

FROM alpine:3.18.4 AS release
WORKDIR /app
COPY --from=build /temp/main .
COPY public public
EXPOSE 8080
ENTRYPOINT ["/app/main"]

# RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-gs-ping