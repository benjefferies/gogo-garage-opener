FROM golang as builder

# Install ARM gcc and build tools
RUN echo "deb http://emdebian.org/tools/debian/ jessie main" >> /etc/apt/sources.list && \
    curl http://emdebian.org/tools/debian/emdebian-toolchain-archive.key | apt-key add - && \
    dpkg --add-architecture armhf && \
    apt-get -y update && \
    apt-get -y install crossbuild-essential-armhf

# The command to use to compile C code.
ENV CC arm-linux-gnueabihf-gcc
# The command to use to compile C++ code.
ENV CXX=arm-linux-gnueabihf-g++
# The OS running on the raspberry pi
ENV GOOS linux
# The CPU architecture of the raspberry pi
ENV GOARCH arm
ENV GOARM 7
# Cgo enables the creation of Go packages that call C code. (required for https://github.com/mattn/go-sqlite3)
ENV CGO_ENABLED 1

ADD gogo-garage-opener /go/src/github.com/benjefferies/gogo-garage-opener

WORKDIR /go/src/github.com/benjefferies/gogo-garage-opener

RUN go get -d -v ./...
RUN go install -v ./...

FROM arm32v7/debian:9-slim

WORKDIR /var/gogo-garage-opener

COPY --from=builder /go/bin/linux_arm/gogo-garage-opener /var/gogo-garage-opener/gogo-garage-opener
COPY gogo-garage-opener/index.html /var/gogo-garage-opener/index.html

CMD [ "./gogo-garage-opener" ]
