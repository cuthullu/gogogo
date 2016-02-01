FROM golang:latest 

ADD ./fakeDataServer.go /go/src/mygo/server/

RUN mkdir /var/lib/dserver
ADD ./data /var/lib/deserver

RUN go get ./src/mygo/server
RUN go install ./src/mygo/server


ENTRYPOINT /go/bin/server -port=8080 -dataDir=/var/lib/deserver

EXPOSE 8080