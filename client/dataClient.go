package main

import (
	"fmt"
	"net"
	"flag"

	"github.com/docker/libchan"
	"github.com/docker/libchan/spdy"
)


type RemoteCommand struct {
    Cmd string
    Args [] string
    OutChan libchan.Sender
    Closer libchan.Receiver

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

    remote := flag.String("remote", "localhost", "Where to find that there server there")
    port := flag.String("port", "8080", "The port to listen on")
    flag.Parse()
    fmt.Println("Config: ", *remote, *port)
    
	client, err := net.Dial("tcp", *remote + ":" + *port)
	check(err)

	p, err := spdy.NewSpdyStreamProvider(client, false)
	check(err)

	transport := spdy.NewTransport(p)
	sender, err := transport.NewSendChannel()
	check(err)

	receiver, remoteSender := libchan.Pipe()
	
	closeReceiver, _ := libchan.Pipe()

	command := &RemoteCommand {
		Cmd : "attach",
		Args : make([]string, 3),
		OutChan : remoteSender,
		Closer : closeReceiver,
	}

	err = sender.Send(command)
	check(err)

	for {
		rLine := &RemoteLine{}
		err = receiver.Receive(rLine)
		fmt.Println(rLine.Line)
	}

	
}