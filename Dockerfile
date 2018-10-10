# first stage
FROM golang:latest as builder
# install glide
RUN go get github.com/Masterminds/glide
# create a working directory
RUN mkdir -p /go/src/drone-digitalocean/bin
ADD . /go/src/drone-digitalocean
WORKDIR /go/src/drone-digitalocean
# install packages
RUN glide install
# build
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o bin/drone-digitalocean .

# second stage
FROM alpine:latest as certs
# get certificates
RUN apk --update add ca-certificates

# final stage
FROM scratch
# copy certificates
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
# copy binary
COPY --from=builder /go/src/drone-digitalocean/bin /app/
# set entrypoint
WORKDIR /app
ENTRYPOINT ["/app/drone-digitalocean"]