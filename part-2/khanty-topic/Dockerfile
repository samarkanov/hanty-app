FROM golang
RUN apt-get update && apt-get -y install vim wget

COPY app /tmp/app
WORKDIR /tmp/app
RUN go get github.com/samarkanov/khanty-app/utils@7405372b8752f2dcee223fe4b6d06d612d379536
RUN go build .

EXPOSE 10001
CMD /tmp/app/topic
