FROM golang:1.12 as builder

RUN curl -L https://github.com/balena-io/qemu/releases/download/v3.0.0%2Bresin/qemu-3.0.0+resin-arm.tar.gz | tar zxvf - -C . && mv qemu-3.0.0+resin-arm/qemu-arm-static .

FROM arm32v7/debian:9-slim

COPY --from=builder /go/src/github.com/benjefferies/gogo-garage-opener/qemu-arm-static /usr/bin

RUN bash -c "echo deb http://ftp.de.debian.org/debian sid main" >> /etc/apt/sources.list
RUN apt-get update
RUN apt-get -y upgrade libc6

WORKDIR /var/gogo-garage-opener

COPY --from=builder /go/bin/linux_arm/gogo-garage-opener /var/gogo-garage-opener/gogo-garage-opener
COPY gogo-garage-opener/index.html /var/gogo-garage-opener/index.html

CMD [ "./gogo-garage-opener" ]
