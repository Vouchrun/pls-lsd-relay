FROM golang:1.22-bullseye AS builder

WORKDIR /app

COPY go.mod go.sum /app/

RUN go mod download

COPY . .

RUN make build

FROM gcr.io/distroless/base-debian12:nonroot

COPY --from=builder /app/build/eth-lsd-relay /app

ENTRYPOINT [ "/app" ] 