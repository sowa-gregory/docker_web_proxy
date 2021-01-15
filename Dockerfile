FROM golang:alpine as build
WORKDIR /app
RUN apk add git binutils
COPY dnsserver.go .
RUN go get -d
RUN go build dnsserver.go
RUN strip dnsserver

FROM alpine:latest
WORKDIR /bin
COPY --from=build /app/dnsserver .
CMD dnsserver
EXPOSE 53

