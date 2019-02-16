FROM arm32v7/golang:1.10 as builder

ADD gogo-garage-opener /go/src/gogo-garage-opener

WORKDIR /go/src/gogo-garage-opener

RUN go get -d -v ./...
RUN go install -v ./...

FROM arm32v7/busybox

COPY --from=builder /go/src/gogo-garage-opener/gogo-garage-opener /var/gogo-garage-opener/gogo-garage-opener

CMD [ "/var/gogo-garage-opener/gogo-garage-opener" ]