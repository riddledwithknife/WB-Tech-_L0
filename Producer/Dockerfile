FROM golang:alpine AS builder-prod

WORKDIR /usr/local/src

COPY go.mod go.sum ./
RUN go mod download

COPY producer/* ./
RUN go build -o ./sender ./sender.go

FROM alpine AS runner-prod

COPY --from=builder-prod /usr/local/src/sender /usr/local/src/order_example.json /

CMD ["/sender"]