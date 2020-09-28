FROM golang
RUN apt-get update && apt-get -y install vim wget

COPY app /tmp/app
COPY config.json /tmp/app
WORKDIR /tmp/app
RUN go build topic.go

EXPOSE 10001
CMD /tmp/app/topic
