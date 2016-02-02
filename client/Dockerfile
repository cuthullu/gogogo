FROM golang:latest 

ADD ./dataClient.go /go/src/gogogo/client/

RUN go get ./src/gogogo/client
RUN go install ./src/gogogo/client

ENTRYPOINT /go/bin/client -remote 172.17.0.2 -port=9900

EXPOSE 8080
