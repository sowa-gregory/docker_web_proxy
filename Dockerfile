FROM golang:alpine as build
WORKDIR /app
RUN apk add git binutils
RUN go get -t github.com/miekg/dns
COPY dnsserver.go .
RUN go build dnsserver.go
RUN strip dnsserver

FROM alpine:latest
WORKDIR /bin
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser
COPY --from=build /app/dnsserver .
ENV PROXY_HOST host
ENTRYPOINT ["dnsserver"]
EXPOSE 53/udp

