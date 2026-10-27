FROM golang:1.27.3 AS builder
WORKDIR /usr/src/pismo-challenge-api
COPY . .
RUN CGO_ENABLED=0 go build -tags timetzdata -o /usr/bin/pismo-challenge-api ./cmd/api-server

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/bin/pismo-challenge-api /usr/bin/
ENTRYPOINT ["pismo-challenge-api"]
