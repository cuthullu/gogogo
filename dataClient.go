package main

import (
	"fmt"
	"net"

	"github.com/docker/libchan"
	"github.com/docker/libchan/spdy"
)


type RemoteCommand struct {
    Cmd string
    Args [] string
    OutChan libchan.Sender

}

type RemoteLine struct {
	Line string
}

func check(e error) {
    if(e != nil) {
        panic(e)
    }
}

func main() {

	client, err := net.Dial("tcp", "127.0.0.1:9323")
	check(err)

	p, err := spdy.NewSpdyStreamProvider(client, false)
	check(err)

	transport := spdy.NewTransport(p)
	sender, err := transport.NewSendChannel()
	check(err)

	receiver, remoteSender := libchan.Pipe()

	command := &RemoteCommand {
		Cmd : "attach",
		Args : make([]string, 3),
		OutChan : remoteSender,
	}

	err = sender.Send(command)
	check(err)

	for {
		rLine := &RemoteLine{}
		err = receiver.Receive(rLine)
		fmt.Println(rLine.Line)
	}
}