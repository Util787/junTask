FROM golang:1.24.3-alpine3.22 AS build

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o user-manager-api ./cmd/main.go

FROM alpine:3.22

COPY --from=build /usr/src/app/user-manager-api .   
COPY --from=build /usr/src/app/.env .   
CMD ["./user-manager-api"]

