package main

import (
    "fmt"
    "io/ioutil"
    "bufio"
    "os"
    "strings"
    "time"
)

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
    tNow = tNow.Add(time.Hour * -7 )

    return tNow.Format("150405")
}

func tickerStreamer(path string, line chan string, finished chan string) {
    
    
    f, err := os.Open(path)
    check(err)
    defer f.Close()
    defer func(){finished <- "Finished reading: " + path}()
    bufr := bufio.NewReader(f)

    //Print header
    dat, _, err := bufr.ReadLine();
    fmt.Println("Reading: " + path)
    fmt.Println(string(dat))
    
    // Fast forward to now
    t := getNow()

    tin := "0"
    for tin < t{
        dat, _, err = bufr.ReadLine();
        if err != nil {
            return
        }
        tin = getTime(string(dat))
    }

    for {
        dat, _, err = bufr.ReadLine();
        if err != nil {
            return
        }
        tin = getTime(string(dat))
        for tin > t {
            time.Sleep(time.Millisecond * 500)
            t = getNow()
        }
        line <- string(dat)

    }
    
}

func readAllTheDatas(dataLine chan string) {
    defer close(dataLine)

    dir := "./data"
    files, _ := ioutil.ReadDir(dir)
    fileFinishes := make(chan string)
    for _, f := range files {
        //f := files[0]
        path := dir + "/" + f.Name()
        go tickerStreamer(path, dataLine, fileFinishes)
    }

    for range files {
        msg := <- fileFinishes
        fmt.Println(msg)
    }
}

func main() {
    dataLine := make(chan string)
    go readAllTheDatas(dataLine)
    for range dataLine {
        line := <- dataLine
        fmt.Println(line)
    }
}