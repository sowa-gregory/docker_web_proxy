FROM golang:alpine as build
WORKDIR /app
RUN apk add git binutils
COPY dnsserver.go .
RUN go get -d
RUN go build dnsserver.go
RUN strip dnsserver

FROM alpine:latest
WORKDIR /bin
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser
COPY --from=build /app/dnsserver .
CMD dnsserver
EXPOSE 53/udp

