package main

import (
    "fmt"
    "io/ioutil"
    "bufio"
    "os"
)

func check(e error) {
    if(e != nil) {
        panic(e)
    }
}

func readFileBuffered(path string, line chan string, finished chan string) {
    
    
    f, err := os.Open(path)
    check(err)
    defer f.Close()
    defer func(){finished <- "Finished reading: " + path}()
    bufr := bufio.NewReader(f)

    dat, _, err := bufr.ReadLine();
    fmt.Println("Reading: " + path)
    fmt.Println(string(dat))
    for {
        dat, _, err = bufr.ReadLine();
        if err != nil {
            return
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
        path := dir + "/" + f.Name()
        go readFileBuffered(path, dataLine, fileFinishes)
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