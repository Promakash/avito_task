FROM golang:1.23.4-alpine AS build

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o avito-shop-service cmd/main.go

FROM alpine AS runner

RUN apk add --no-cache curl

WORKDIR /app

COPY --from=build /build/avito-shop-service ./
COPY --from=build /build/config ./config/

EXPOSE 8080

CMD ["./avito-shop-service"]