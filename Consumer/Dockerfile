FROM golang:alpine AS builder-cons

WORKDIR /usr/local/src

COPY go.mod go.sum ./
RUN go mod download

COPY consumer/* ./
RUN go build -o ./service ./service.go

FROM alpine AS runner-cons

COPY --from=builder-cons /usr/local/src/service /usr/local/src/scheme.json /

CMD ["/service"]