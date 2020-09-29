FROM golang
RUN apt-get update && apt-get -y install vim wget

COPY app /tmp/app
COPY config.json /tmp/app/config.json
WORKDIR /tmp/app
RUN go build .

EXPOSE 10002
CMD /tmp/app/appconfig
