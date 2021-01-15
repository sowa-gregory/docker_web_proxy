FROM golang as build
WORKDIR /app
COPY dnsserver.go .
RUN go build show_flag.go
RUN strip show_flag

FROM alpine:latest
WORKDIR /bin
COPY --from=build /app/show_flag .
