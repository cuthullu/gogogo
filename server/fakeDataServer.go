package main

import (
    "fmt"
    "io/ioutil"
    "bufio"
    "os"
    "strings"
    "time"
    "net"
    "flag"

    "github.com/docker/libchan"
    "github.com/docker/libchan/spdy"
)

type RemoteCommand struct {
    Cmd string
    Args []string
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

func getTicker(line string) string {
    return strings.Split(line, ",")[0]
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

func readAllTheDatas(dir string, outChan chan string) {
    defer close(outChan)
    fileFinishes := make(chan string)

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

func listen(listener net.Listener, clients *[]libchan.Sender) {
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

                *clients = append(*clients, command.OutChan)

            }
        }()
    }
}


func addOutListener(outs *[]int) {
    *outs = append(*outs, 5)
}   

func main() {
    outChan := make(chan string)

    dataDir := flag.String("dataDir", "data", "The directory what which has the datas in it")
    hostAdd := flag.String("hostAdd", "localhost", "The host address to listen on")
    port := flag.String("port", "9323", "The port to listen on")
    flag.Parse()
    fmt.Println("Config: ", *dataDir, *hostAdd, *port)
    

    go readAllTheDatas(*dataDir, outChan)

    var listener net.Listener
    var err error

    var clients []libchan.Sender


    listener, err = net.Listen("tcp", *hostAdd + ":" + *port)

    check(err)
    go listen(listener, &clients)

    for range outChan {
        line := <- outChan
        rLine := &RemoteLine{
            Line : line,
        }
        for _, client := range clients {
            client.Send(rLine)
        }


    }
}