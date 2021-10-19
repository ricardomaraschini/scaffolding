FROM docker.io/library/golang:1.16 AS builder
WORKDIR /go/src/app
COPY . .
RUN make

FROM docker.io/library/fedora:latest
COPY --from=builder /go/src/app/_output/bin/app /usr/local/bin/app
EXPOSE 8080
ENTRYPOINT [ "/usr/local/bin/app" ]
CMD [ "serve" ]
