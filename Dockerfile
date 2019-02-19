FROM golang as builder

RUN echo "deb http://emdebian.org/tools/debian/ jessie main" >> /etc/apt/sources.list && \
    curl http://emdebian.org/tools/debian/emdebian-toolchain-archive.key | apt-key add - && \
    dpkg --add-architecture armhf && \
    apt-get -y update && \
    apt-get -y install crossbuild-essential-armhf

ENV CC arm-linux-gnueabihf-gcc
ENV CXX=arm-linux-gnueabihf-g++
ENV GOOS linux
ENV GOARCH arm
ENV GOARM 7
ENV CGO_ENABLED 1

ADD gogo-garage-opener /go/src/github.com/benjefferies/gogo-garage-opener

WORKDIR /go/src/github.com/benjefferies/gogo-garage-opener

RUN go get -d -v ./...
RUN go install -v ./...

FROM arm32v7/debian:9-slim

COPY --from=builder /go/bin/linux_arm/gogo-garage-opener /var/gogo-garage-opener/gogo-garage-opener

CMD [ "/var/gogo-garage-opener/gogo-garage-opener" ]
