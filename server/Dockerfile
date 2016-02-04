FROM golang:latest 

ADD ./fakeDataServer.go /go/src/gogogo/server/

RUN mkdir /var/lib/dserver
ADD ./data /var/lib/deserver

RUN go get ./src/gogogo/server
RUN go install ./src/gogogo/server

#ENTRYPOINT /go/bin/server -hostAdd 172.17.0.2 -port=9900 -dataDir=/var/lib/deserver

#EXPOSE 8080
