package main

import (
    "fmt"
    "io/ioutil"
    "bufio"
    "os"
    "strings"
    "time"
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

func getTime(line string) string{
    return strings.Split(line, ",")[3]
}

func getNow() string {
    tNow := time.Now()

    //Cuz I work harder than the markets
    //tNow = tNow.Add(time.Hour * -7 )

    return tNow.Format("150405")
}

func tickerStreamer(path string, outChan chan string, finished chan string) {
    
    f, err := os.Open(path)
    check(err)
    defer f.Close()
    defer func(){finished <- "Finished reading: " + path}()
    bufr := bufio.NewReader(f)

    //Throw away header
    dat, _, err := bufr.ReadLine();
    
    tNow := getNow()
    tin := "0"

    //Fast forward input to current time
    for tin < tNow{
        dat, _, err = bufr.ReadLine();
        if err != nil {
            return
        }
        tin = getTime(string(dat))
    }

    // Print input at current second. Sleep if input is after current second
    for {
        dat, _, err = bufr.ReadLine();
        if err != nil {
            return
        }
        tin = getTime(string(dat))
        for tin > tNow {
            time.Sleep(time.Millisecond * 500)
            tNow = getNow()
        }
        outChan <- string(dat)

    }
    
}

func readAllTheDatas(outChan chan string) {
    defer close(outChan)
    fileFinishes := make(chan string)

    dir := "./data"
    files, _ := ioutil.ReadDir(dir)
    
    for _, f := range files {
        path := dir + "/" + f.Name()
        go tickerStreamer(path, outChan, fileFinishes)
    }

    for range files {
        msg := <- fileFinishes
        fmt.Println(msg)
    }
}

func listen(listener net.Listener, outs *[]libchan.Sender) {
    for {
        c, err := listener.Accept()
        check(err)

        p, err := spdy.NewSpdyStreamProvider(c, true)
        check(err)

        t := spdy.NewTransport(p)

        go func() {
            for {
                receiver, err := t.WaitReceiveChannel()
                check(err)
                command := &RemoteCommand{}
                err = receiver.Receive(command)
                check(err)

                *outs = append(*outs, command.OutChan)

            }
        }()
    }
}


func addOutListener(outs *[]int) {
    *outs = append(*outs, 5)
}   

func main() {
    outChan := make(chan string)
    go readAllTheDatas(outChan)

    var listener net.Listener
    var err error

    var outs []libchan.Sender


    listener, err = net.Listen("tcp", "localhost:9323")

    check(err)
    go listen(listener, &outs)

    for range outChan {
        line := <- outChan
        fmt.Println(line)
        rLine := &RemoteLine{
            Line : line,
        }
        for _, out := range outs {
            out.Send(rLine)
        }

    }
}