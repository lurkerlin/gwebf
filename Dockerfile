FROM golang:latest

WORKDIR /build
COPY . .
RUN go build -o app
WORKDIR /dist
RUN cp /build/app .
EXPOSE 8080
CMD ["/dist/app"]