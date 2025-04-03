FROM golang:1.21-alpine as build
WORKDIR /go/src/github.com/GeoNet/fits/dapper
COPY ./cmd/dapper-api ./cmd/dapper-api
COPY ./dapperlib ./dapperlib
COPY ./internal ./internal
COPY ./vendor ./vendor
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags "-X main.Prefix=/usr/local -extldflags -static" -installsuffix cgo -o dapper-api cmd/dapper-api/*.go
RUN apk --update add ca-certificates
FROM scratch
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /go/src/github.com/GeoNet/fits/dapper/dapper-api ./dapper-api
COPY ./cmd/dapper-api/assets ./assets/
CMD ["./dapper-api"]
