package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

func main() {
	addr := net.JoinHostPort("localhost", "9000")
	tcpClient(addr)
}

func tcpClient(addr string) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		Quitf("Could not open connection to server: %#v", err)
	}

	var (
		buf = make([]byte, 1024)
		r   = bufio.NewReader(conn)
		w   = bufio.NewWriter(conn)
	)

	for {
		n, err := r.Read(buf)
		data := string(buf[:n])

		switch err {
		case io.EOF:
			Quitf("Connection broke")
		case nil:
			log.Printf("received data: %s", data)
			if strings.Contains(data, "PONG") {
				time.Sleep(time.Second)
				sendPing(w)
			}
		default:
			Quitf("Receiving data failing: %#v", err)
		}
	}
}

func sendPing(w *bufio.Writer) {
	_, err := w.WriteString("PING" + "\r\n")
	if err != nil {
		Quitf("Could not send PING: %#v", err)
	}

	err = w.Flush()
	if err != nil {
		Quitf("Could not send PING: %#v", err)
	}
}

func Quitf(format string, a ...interface{}) {
	fmt.Printf(format, a...)
	os.Exit(1)
}
