package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

var quitSemaphore chan bool

func main() {
	var laddr = net.TCPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Zone: "",
	}

	raddr := net.TCPAddr{IP: net.ParseIP("127.0.0.1"),
		Port: 3384,
		Zone: ""}

	tcp, err := net.DialTCP("tcp4", &laddr, &raddr)
	if err != nil {
		panic(err)
	}
	defer tcp.Close()

	go write(tcp)

	<-quitSemaphore

}

func write(conn *net.TCPConn) {
	reader := bufio.NewReader(os.Stdin)
	newReader := bufio.NewReader(conn)

	defer func() {
		quitSemaphore <- false
	}()

	for {
		stdin, err := reader.ReadString(byte('\n'))
		if err != nil {
			log.Println(err)
			return
		}

		_, err = conn.Write([]byte(stdin))
		if err != nil {
			fmt.Println(err)
			return
		}

		readString, err := newReader.ReadString(byte('\n'))
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(readString)

	}
}
