FROM golang:1.25 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o app_pull_assign ./app_pull_assign

FROM alpine:latest

WORKDIR /app

COPY --from=build /app/app_pull_assign .

RUN chmod +x ./app_pull_assign

WORKDIR /root/

CMD ["/app/app_pull_assign"]

EXPOSE 8080
