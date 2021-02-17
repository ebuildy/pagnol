FROM golang:1.14.3-alpine AS builder

RUN apk --no-cache add tzdata zip ca-certificates && \
     update-ca-certificates

WORKDIR /src

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build \
    -ldflags='-w -s -extldflags "-static"' -a \
    -o /go/bin/pagnol

FROM scratch

ENV ZONEINFO /zoneinfo.zip

COPY --from=builder /go/bin/pagnol /go/bin/pagnol
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["/go/bin/pagnol"]