FROM golang
RUN apt-get update && apt-get -y install vim wget

COPY app /tmp/app
COPY config.json /tmp/app/config.json
WORKDIR /tmp/app
RUN go get github.com/samarkanov/khanty-app/utils@c4de44d1224690a78f2c4d80527c014ed600cafb
RUN go build .

EXPOSE 10002
CMD /tmp/app/appconfig
