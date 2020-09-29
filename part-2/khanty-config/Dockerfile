FROM golang
RUN apt-get update && apt-get -y install vim wget

COPY app /tmp/app
COPY config.json /tmp/app/config.json
WORKDIR /tmp/app
RUN go get github.com/samarkanov/khanty-app/utils@bd207ac5fddec7235e73dea916fc35211d9aa982
RUN go build .

EXPOSE 10002
CMD /tmp/app/appconfig
