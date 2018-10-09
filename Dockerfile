FROM golang:alpine as builder
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN curl https://glide.sh/get | sh
RUN glide install
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main .

FROM scratch
COPY --from=builder /build/main /app/
WORKDIR /app
CMD ["./main"]