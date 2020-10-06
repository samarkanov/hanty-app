FROM golang
RUN apt-get update && apt-get -y install vim wget

ENV KHANTY_CONFIG_FILE /mount/.config/config.json
COPY app /tmp/app
WORKDIR /tmp/app
RUN go get github.com/samarkanov/khanty-app/utils
RUN go build .
