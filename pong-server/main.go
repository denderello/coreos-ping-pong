package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 3 {
		Quitf("You need to pass host and port as parameters")
	}

	addr := net.JoinHostPort(os.Args[1], os.Args[2])
	tcpServer(addr)
}

func tcpServer(addr string) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		Quitf("Could not start server: %#v", err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			Quitf("Error handling connection: %#v", err)
		}
		go handler(conn)
	}
}

func handler(conn net.Conn) {
	defer conn.Close()

	var (
		buf = make([]byte, 1024)
		r   = bufio.NewReader(conn)
		w   = bufio.NewWriter(conn)
	)

	writeMessage(w, "Welcome to the PING/PONG server")

	for {
		n, err := r.Read(buf)
		data := string(buf[:n])

		switch err {
		case io.EOF:
			Quitf("Connection broke")
		case nil:
			log.Printf("received data: %s", data)
			if strings.Contains(data, "PING") {
				sendPong(w)
			}
		default:
			Quitf("Receiving data failing: %#v", err)
		}
	}
}

func sendPong(w *bufio.Writer) {
	writeMessage(w, "PONG")
}

func writeMessage(w *bufio.Writer, m string) {
	_, err := w.WriteString(m + "\r\n")
	if err != nil {
		Quitf("Could not send PONG: %#v", err)
	}

	err = w.Flush()
	if err != nil {
		Quitf("Could not send PONG: %#v", err)
	}
}

func Quitf(format string, a ...interface{}) {
	fmt.Printf(format, a...)
	os.Exit(1)
}
