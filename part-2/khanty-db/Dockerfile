FROM golang
RUN apt-get update && apt-get -y install vim wget

COPY app /tmp/app
WORKDIR /tmp/app
RUN go build .

EXPOSE 10003
CMD /tmp/app/db
