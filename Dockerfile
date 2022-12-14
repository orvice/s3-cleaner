FROM golang:1.19 as builder

RUN mkdir /go/src/app
WORKDIR /go/src/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLE=0 go build -o /go/bin/app ./cmd/main.go

FROM gcr.io/distroless/base-debian10

LABEL org.opencontainers.image.description email:git@orx.me


COPY --from=builder /go/bin/app /app