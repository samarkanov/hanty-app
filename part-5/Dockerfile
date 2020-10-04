FROM golang
RUN apt-get update && apt-get -y install vim wget

COPY app /tmp/app
WORKDIR /tmp/app
RUN go get github.com/samarkanov/khanty-app/utils@7d6d1201c47ca805fd5c84de698aea230adcadf5
RUN go build .

EXPOSE 10004
CMD /tmp/app/khanty-master
