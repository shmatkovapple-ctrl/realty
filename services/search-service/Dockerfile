FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.work go.work.sum ./

COPY go.mod ./go.mod

COPY services/user-service/go.mod     services/user-service/go.mod
COPY services/listing-service/go.mod  services/listing-service/go.mod
COPY services/deal-service/go.mod     services/deal-service/go.mod
COPY services/search-service/go.mod   services/search-service/go.mod
COPY services/notification-service/go.mod services/notification-service/go.mod
COPY services/api-gateway/go.mod      services/api-gateway/go.mod

COPY services/user-service/go.sum     services/user-service/go.sum
COPY services/listing-service/go.sum  services/listing-service/go.sum
COPY services/deal-service/go.sum     services/deal-service/go.sum
COPY services/search-service/go.sum   services/search-service/go.sum
COPY services/notification-service/go.sum services/notification-service/go.sum
COPY services/api-gateway/go.sum      services/api-gateway/go.sum

RUN go work sync

COPY api/ api/
COPY pkg/ pkg/
COPY services/ services/

RUN go build -o /app/service realty/services/search-service/cmd

FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/service .

ENV TZ=Europe/Moscow

EXPOSE 50054

CMD ["./service"]