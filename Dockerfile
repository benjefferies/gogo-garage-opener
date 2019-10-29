FROM arm32v7/debian:9-slim

RUN apt-get update
RUN apt-get -y upgrade libc6

WORKDIR /var/gogo-garage-opener

CMD [ "./gogo-garage-opener" ]
