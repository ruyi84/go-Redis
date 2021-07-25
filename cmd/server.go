package main

import "C"

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

var m = map[string]string{}

func main() {
	addr := &net.TCPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: 3384,
		Zone: "",
	}

	tcp, err := net.ListenTCP("tcp4", addr)
	if err != nil {
		panic(err)
	}

	defer tcp.Close()

	for {
		acceptTCP, err := tcp.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		go tcpRead(acceptTCP)
	}

}

func tcpRead(conn net.Conn) {
	reader := bufio.NewReader(conn)

	for {
		readString, err := reader.ReadString(byte('\n'))
		if err != nil {
			fmt.Println("close")
			return
		}

		trimRight := strings.TrimRight(readString, "\n")

		split := strings.Split(trimRight, " ")
		fmt.Println(split)
		switch split[0] {
		case "Ping":
			conn.Write([]byte("Pong \n"))
			continue
		case "Set":
			if len(split) != 3 {
				conn.Write([]byte("set parsing is wrong \n"))
				continue
			}
			m[split[1]] = split[2]
			conn.Write([]byte("Ok \n"))
			continue
		case "Del":
			if len(split) != 2 {
				conn.Write([]byte("set parsing is wrong \n"))
				continue
			}
			delete(m, split[1])
			conn.Write([]byte("Ok \n"))
			continue
		case "Get":
			if len(split) != 2 {
				conn.Write([]byte("get parsing is wrong \n"))
				continue
			}
			result := fmt.Sprintf("%s ", m[split[1]])
			conn.Write([]byte(result + "\n"))
			continue
		}
		conn.Write([]byte("wrong \n"))
	}
}

func Set() {
	fmt.Println("adsf")
}
