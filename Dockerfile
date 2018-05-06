FROM golang:latest

ENV CC arm-linux-gnueabihf-gcc
ENV GOOS linux
ENV GOARCH arm
ENV GOARM 6
ENV CGO_ENABLED 1

RUN echo "deb http://emdebian.org/tools/debian/ jessie main" >> /etc/apt/sources.list && \
    curl http://emdebian.org/tools/debian/emdebian-toolchain-archive.key | apt-key add - && \
    dpkg --add-architecture armhf && \
    apt-get -y update && \
    apt-get -y install crossbuild-essential-armhf

CMD go get -v && go build -v -o gogo-garage-opener
