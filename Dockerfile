FROM golang:1.19 as builder

RUN mkdir /go/src/app
WORKDIR /go/src/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /go/bin/app ./cmd/main.go

FROM gcr.io/distroless/base-debian10
COPY --from=builder /go/bin/app /app