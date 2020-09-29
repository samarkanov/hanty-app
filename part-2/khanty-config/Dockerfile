FROM golang
RUN apt-get update && apt-get -y install vim wget

COPY app /tmp/app
COPY config.json /tmp/app/config.json
WORKDIR /tmp/app
RUN go get github.com/samarkanov/khanty-app/utils@7d6d1201c47ca805fd5c84de698aea230adcadf5
RUN go build .

EXPOSE 10002
CMD /tmp/app/appconfig
