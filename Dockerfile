FROM golang:latest as builder

ENV CC arm-linux-gnueabihf-gcc
ENV GOOS linux
ENV GOARCH arm
ENV GOARM 6
ENV CGO_ENABLED 1

ADD gogo-garage-opener /go/src/gogo-garage-opener

WORKDIR /go/src/gogo-garage-opener

RUN echo "deb http://emdebian.org/tools/debian/ jessie main" >> /etc/apt/sources.list && \
    curl http://emdebian.org/tools/debian/emdebian-toolchain-archive.key | apt-key add - && \
    dpkg --add-architecture armhf && \
    apt-get -y update && \
    apt-get -y install crossbuild-essential-armhf && \
    go get -v && go build -v -o gogo-garage-opener

FROM scratch

COPY --from=builder /go/src/gogo-garage-opener/gogo-garage-opener /var/gogo-garage-opener/gogo-garage-opener

CMD [ "/var/gogo-garage-opener/gogo-garage-opener" ]