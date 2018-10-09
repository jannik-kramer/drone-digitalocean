# first stage
FROM golang:latest as builder

# install glide
RUN go get github.com/Masterminds/glide

# create a working directory
RUN mkdir -p /go/src/drone-digitalocean
ADD . /go/src/drone-digitalocean
WORKDIR /go/src/drone-digitalocean

# install packages
RUN glide install

# build
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main .

# second stage
FROM scratch

# copy binary
COPY --from=builder /go/src/drone-digitalocean /app/

# set entrypoint
WORKDIR /app
CMD ["./app/main"]