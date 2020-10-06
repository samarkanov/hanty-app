FROM golang
RUN apt-get update && apt-get -y install vim wget

COPY app /tmp/app
WORKDIR /tmp/app
RUN go get github.com/samarkanov/khanty-app/utils
RUN go build .
