FROM golang:alpine as builder
RUN apk add --no-cache ca-certificates git

WORKDIR /src
COPY ./ ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w"

FROM scratch as runtime
COPY --from=builder /src/covid-decoder ./
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/covid-decoder"]